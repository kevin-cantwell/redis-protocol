package resp

import (
	"bufio"
	"fmt"
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

func (w *Writer) WriteError(value string) error {
	defer w.out.Flush()

	if err := w.out.WriteByte(errorPrefix); err != nil {
		return err
	}
	if _, err := w.out.WriteString(value); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteSimpleString(value string) error {
	defer w.out.Flush()

	if err := w.out.WriteByte(simpleStringPrefix); err != nil {
		return err
	}
	if _, err := w.out.WriteString(value); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBulkString(value *string) error {
	defer w.out.Flush()

	if err := w.out.WriteByte(bulkStringPrefix); err != nil {
		return err
	}
	if value == nil {
		if _, err := w.out.WriteString("-1\r\n"); err != nil {
			return err
		}
		return nil
	}
	if _, err := w.out.WriteString(fmt.Sprint(len(*value))); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	if _, err := w.out.WriteString(*value); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteInteger(value int64) error {
	defer w.out.Flush()

	if err := w.out.WriteByte(integerPrefix); err != nil {
		return err
	}
	if _, err := w.out.WriteString(fmt.Sprint(value)); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteArray(array ...Data) error {
	defer w.out.Flush()

	if err := w.out.WriteByte(arrayPrefix); err != nil {
		return err
	}
	if _, err := w.out.WriteString(fmt.Sprint(len(array))); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	for _, data := range array {
		if err := w.WriteData(data); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) WriteData(data Data) error {
	switch d := data.(type) {
	case SimpleString:
		return w.WriteSimpleString(d.Val)
	case BulkString:
		return w.WriteBulkString(d.Val)
	case Error:
		return w.WriteError(d.Val)
	case Integer:
		return w.WriteInteger(d.Val)
	case Array:
		return w.WriteArray(d.Val...)
	default:
		return fmt.Errorf("resp: unknown data type %T", d)
	}
}
