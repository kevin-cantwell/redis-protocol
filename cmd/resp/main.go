package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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
	app.Action = func(c *cli.Context) {
		var err error
		var reader io.Reader = os.Stdin

		// If a filename was passed in as an argument, read from the file instead of stdin
		if len(c.Args()) > 0 {
			reader, err = os.Open(c.Args()[0])
			if err != nil {
				exit(err.Error(), 1)
			}
		}

		scanner := bufio.NewScanner(reader)
		respWriter := NewRESPWriter(os.Stdout)

		for scanner.Scan() {
			line := scanner.Text()

			args := strings.Split(line, " ")
			if err := respWriter.WriteCommand(args...); err != nil {
				exit(err.Error(), 1)
			}
		}
		if err := scanner.Err(); err != nil {
			exit(err.Error(), 1)
		}
	}

	app.Run(os.Args)
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
