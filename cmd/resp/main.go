package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

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

func skipEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}

func decodeRESP(ctx *cli.Context, r io.Reader) error {
	reader := resp.NewReader(r)
	for {
		data, err := reader.ReadData()
		if err != nil {
			return skipEOF(err)
		}
		var output string
		if ctx.Bool("raw") {
			output = data.Raw()
		} else {
			output = data.Human()
		}
		if _, err := os.Stdout.WriteString(output + "\n"); err != nil {
			return err
		}
	}
}

func encodeRESP(ctx *cli.Context, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	respWriter := resp.NewWriter(os.Stdout)

	for scanner.Scan() {
		var array resp.Array
		fields, err := parseFields(scanner.Bytes())
		if err != nil {
			return err
		}
		for _, field := range fields {
			array = append(array, resp.BulkString(field))
		}
		if err := respWriter.WriteData(array); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func parseFields(line []byte) ([][]byte, error) {
	var field []byte
	var fields [][]byte
	for i := 0; i < len(line); i++ {
		b := line[i]
		switch b {
		case ' ': // Treat 1 or more spaces as a delim
			if len(field) > 0 {
				fields = append(fields, field)
				field = nil
			}
		case '"': // Treat double quoted strings as a single field and unescape chars like tabs or newlines
			j := bytes.Index(line[i+1:], []byte{b})
			if j < 0 {
				field = append(field, b)
				continue
			}
			// Append everything including the terminating quote
			unquoted, err := strconv.Unquote(string(line[i : i+j+2]))
			if err != nil {
				return nil, err
			}
			field = append(field, []byte(unquoted)...)
			i += j + 1
		case '\'': // Treat single quotes as literal strings
			j := bytes.Index(line[i+1:], []byte{b})
			if j < 0 {
				field = append(field, b)
				continue
			}
			// Append everything including the terminating quote
			field = append(field, line[i+1:i+j+1]...)
			i += j + 1
		default:
			field = append(field, b)
		}
	}
	if len(field) > 0 {
		fields = append(fields, field)
	}
	return fields, nil
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
