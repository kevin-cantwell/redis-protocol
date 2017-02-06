package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

var (
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
)

type RESPWriter struct {
	*bufio.Writer
}

func NewRESPWriter(writer io.Writer) *RESPWriter {
	return &RESPWriter{
		Writer: bufio.NewWriter(writer),
	}
}

func (w *RESPWriter) WriteCommand(args ...string) (err error) {
	// Write the array prefix and the number of arguments in the array.
	w.Write(arrayPrefixSlice)
	w.WriteString(strconv.Itoa(len(args)))
	w.Write(lineEndingSlice)

	// Write a bulk string for each argument.
	for _, arg := range args {
		w.Write(bulkStringPrefixSlice)
		w.WriteString(strconv.Itoa(len(arg)))
		w.Write(lineEndingSlice)
		w.WriteString(arg)
		w.Write(lineEndingSlice)
	}

	return w.Flush()
}

func main() {
	app := cli.NewApp()
	app.Name = "redis-protocol"
	app.Usage = "A parser that converts a file of redis commands to the redis protocol."
	app.Action = func(c *cli.Context) error {
		var err error
		var reader io.Reader = os.Stdin

		// If a filename was passed in as an argument, read from the file instead of stdin
		if len(c.Args()) > 0 {
			reader, err = os.Open(c.Args()[0])
			if err != nil {
				return err
			}
		}

		scanner := bufio.NewScanner(reader)
		respWriter := NewRESPWriter(os.Stdout)

		for scanner.Scan() {
			var args []string
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
					args = append(args, string(arg))
					arg = nil
				case ' ':
					args = append(args, string(arg))
					arg = nil
				default:
					arg = append(arg, b)
				}
			}
			if arg != nil {
				args = append(args, string(arg))
			}

			if err := respWriter.WriteCommand(args...); err != nil {
				return err
			}
		}
		return scanner.Err()
	}

	if err := app.Run(os.Args); err != nil {
		exit(err.Error(), 1)
	}
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
