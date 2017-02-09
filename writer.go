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

func (w *Writer) WriteBulkString(value string) error {
	if err := w.out.WriteByte(bulkStringPrefix); err != nil {
		return err
	}
	if _, err := w.out.WriteString(fmt.Sprint(len(value))); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
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

func (w *Writer) WriteInteger(value int64) error {
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
	if err := w.out.WriteByte(arrayPrefix); err != nil {
		return err
	}
	if _, err := w.out.WriteString(fmt.Sprint(len(array))); err != nil {
		return err
	}
	if _, err := w.out.WriteString("\r\n"); err != nil {
		return err
	}
	for range array {
		// switch v := data.(type) {
		// case string, []byte:
		// 	if err := w.WriteBulkString(string(v)); err != nil {
		// 		return err
		// 	}
		// case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		// 	if err := w.WriteBulkString(string(v)); err != nil {
		// 		return err
		// 	}
		// }
	}
	return nil
}

// func (w *Writer) WriteData(data Data) error {
//   switch d := data.(type) {
//   case *SimpleString:
//   }
// }
