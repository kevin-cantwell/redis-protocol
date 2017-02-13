package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/kevin-cantwell/resp"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "resp"
	app.Usage = "A redis protocol codex."
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "decode, d",
			Usage: "Decodes redis protocol. Default is to encode.",
		},
		cli.BoolFlag{
			Name:  "raw, r",
			Usage: "Decodes redis protocol into a raw format. Default is human-readable.",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		var err error
		var reader io.Reader = os.Stdin

		// If a filename was passed in as an argument, read from the file instead of stdin
		if len(ctx.Args()) > 0 {
			reader, err = os.Open(ctx.Args()[0])
			if err != nil {
				return err
			}
		}

		if ctx.Bool("decode") {
			return decodeRESP(ctx, reader)
		} else {
			return encodeRESP(ctx, reader)
		}
	}

	if err := app.Run(os.Args); err != nil {
		exit(err.Error(), 1)
	}
}

func decodeRESP(ctx *cli.Context, r io.Reader) error {
	reader := resp.NewReader(r)
	for {
		data, err := reader.ReadData()
		if err != nil {
			return skipEOF(err)
		}
		if err := decodeData(ctx, data); err != nil {
			return skipEOF(err)
		}
	}
}

func skipEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}

func decodeData(ctx *cli.Context, data resp.Data) error {
	var d string
	if ctx.Bool("raw") {
		d = data.Raw()
	} else {
		d = data.Human()
	}
	_, err := os.Stdout.WriteString(fmt.Sprintf("%s\n", d))
	return err
}

func encodeRESP(ctx *cli.Context, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	respWriter := resp.NewWriter(os.Stdout)

	for scanner.Scan() {
		var array resp.Array
		var arg []byte
		scanned := scanner.Bytes()
		for i := 0; i < len(scanned); i++ {
			b := scanned[i]
			switch b {
			case '\'', '"':
				// Loop through until we find a terminating quote
				arg = append(arg, b)
				for i++; i < len(scanned); i++ {
					c := scanned[i]
					if c == b {
						arg = arg[1:] // If the quote is terminated, strip the leading quote char
						break
					}
					arg = append(arg, c)
				}
				array = append(array, resp.BulkString(arg))
				arg = nil
			case ' ':
				array = append(array, resp.BulkString(arg))
				arg = nil
			default:
				arg = append(arg, b)
			}
		}
		if arg != nil {
			array = append(array, resp.BulkString(arg))
		}

		if err := respWriter.WriteData(array); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
