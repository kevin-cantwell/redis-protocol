package resp

import (
	"bufio"
	"io"
)

type Writer struct {
	out *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		out: bufio.NewWriter(w),
	}
}

func (w *Writer) WriteData(data Data) error {
	defer w.out.Flush()
	if data == nil {
		_, err := w.out.WriteString(BulkString(nil).Protocol())
		return err
	}
	_, err := w.out.WriteString(data.Protocol())
	return err
}

func (w *Writer) WriteError(value string) error {
	return w.WriteData(Error(value))
}

func (w *Writer) WriteSimpleString(value string) error {
	return w.WriteData(SimpleString(value))
}

func (w *Writer) WriteBulkString(value string) error {
	return w.WriteData(BulkString(value))
}

func (w *Writer) WriteNil() error {
	return w.WriteData(nil)
}

func (w *Writer) WriteInteger(value int64) error {
	return w.WriteData(Integer(value))
}

func (w *Writer) WriteArray(array ...Data) error {
	return w.WriteData(Array(array))
}
