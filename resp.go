package resp

import (
	"fmt"
	"strconv"
	"strings"
)

type Data interface {
	// Protocol returns the full Redis protocol representation including prefix and trailing CRLF.
	Protocol() string
	// Raw is the raw string representation. Nil bulk strings return as empty strings.
	Raw() string
	// Human is a human-readable representation.
	Human() string
	// Quote is a single-line, quoted representation. Array elements are quoted individually.
	Quote() string
}

type Error string

func (d Error) Protocol() string {
	return fmt.Sprintf("-%s\r\n", d)
}

func (d Error) Raw() string {
	return string(d)
}

func (d Error) Human() string {
	return fmt.Sprintf("(error) %s", d)
}

func (d Error) Quote() string {
	return strconv.Quote(d.Raw())
}

type SimpleString string

// should we validate newline chars?
func (d SimpleString) Protocol() string {
	return fmt.Sprintf("+%s\r\n", d)
}

func (d SimpleString) Raw() string {
	return string(d)
}

func (d SimpleString) Human() string {
	return string(d)
}

func (d SimpleString) Quote() string {
	return strconv.Quote(d.Raw())
}

type BulkString []byte

func (d BulkString) Protocol() string {
	if d == nil {
		return "$-1\r\n"
	}
	s := string(d)
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func (d BulkString) Raw() string {
	return string(d)
}

func (d BulkString) Human() string {
	if d == nil {
		return "(nil)"
	}
	return fmt.Sprintf("%q", d)
}

func (d BulkString) Quote() string {
	return strconv.Quote(d.Raw())
}

type Integer int64

func (d Integer) Protocol() string {
	return fmt.Sprintf(":%d\r\n", d)
}

func (d Integer) Raw() string {
	return fmt.Sprintf("%d", d)
}

func (d Integer) Human() string {
	// quoted
	return fmt.Sprintf("%d", d)
}

func (d Integer) Quote() string {
	return strconv.Quote(d.Raw())
}

type Array []Data

func (d Array) Protocol() string {
	s := fmt.Sprintf("*%d\r\n", len(d))
	for _, data := range d {
		s += data.Protocol()
	}
	return s
}

func (d Array) Raw() string {
	var s []string
	for _, data := range d {
		s = append(s, data.Raw())
	}
	return strings.Join(s, "\n")
}

func (d Array) Human() string {
	// quoted
	var s []string
	for i, data := range d {
		s = append(s, fmt.Sprintf("%d) %s", i+1, data.Human()))
	}
	return strings.Join(s, "\n")
}

func (d Array) Quote() string {
	quotes := make([]string, len(d))
	for i, e := range d {
		quotes[i] = e.Quote()
	}
	return strings.Join(quotes, " ")
}
