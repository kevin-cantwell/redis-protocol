package resp

import "fmt"

const (
	simpleStringPrefix byte = '+'
	errorPrefix             = '-'
	integerPrefix           = ':'
	bulkStringPrefix        = '$'
	arrayPrefix             = '*'
)

type Data interface {
	Protocol() string // Only this package may implement
}

type Error string

func (d Error) Protocol() string {
	return fmt.Sprintf("-%s\r\n", d)
}

type SimpleString string

// should we validate newline chars?
func (d SimpleString) Protocol() string {
	return fmt.Sprintf("+%s\r\n", d)
}

// BulkString supports null values so only the pointer type implements the Data interface
type BulkString []byte

func (d BulkString) Protocol() string {
	if d == nil {
		return "$-1\r\n"
	}
	s := string(d)
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

type Integer int64

func (d Integer) Protocol() string {
	return fmt.Sprintf(":%d\r\n", d)
}

type Array []Data

func (d Array) Protocol() string {
	s := fmt.Sprintf("*%d\r\n", len(d))
	for _, data := range d {
		s += data.Protocol()
	}
	return s
}
