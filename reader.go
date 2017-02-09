package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	src *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		src: bufio.NewReader(r),
	}
}

func (r *Reader) ReadData() (Data, error) {
	for {
		b, err := r.src.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		switch b {
		case simpleStringPrefix:
			return r.readSimpleString()
		case errorPrefix:
			return r.readError()
		case integerPrefix:
			return r.readInteger()
		case bulkStringPrefix:
			return r.readBulkString()
		case arrayPrefix:
			return r.readArray()
		default:
			return nil, fmt.Errorf("resp: unknown data type %q", b)
		}
	}
}

func (r *Reader) readSimpleString() (Data, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	return SimpleString{Val: string(line)}, nil
}

func (r *Reader) readError() (Data, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	return Error{Val: string(line)}, nil
}

func (r *Reader) readInteger() (Data, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	i, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return nil, err
	}
	return Integer{Val: i}, nil
}

func (r *Reader) readBulkString() (Data, error) {
	line1, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(string(line1))
	if err != nil {
		return nil, err
	}
	if size == -1 {
		return BulkString{Val: nil}, nil
	}
	line2 := make([]byte, size)
	for i := 0; i < size+2; i++ { // The extra 2 bytes is for the line terminator
		b, err := r.src.ReadByte()
		if err != nil {
			return nil, err
		}
		if i < size {
			line2[i] = b
		}
		if (i == size && b != '\r') || (i == size+1 && b != '\n') {
			return nil, fmt.Errorf("resp: invalid bulk string terminator %q", b)
		}
	}
	value := string(line2)
	return BulkString{Val: &value}, nil
}

func (r *Reader) readArray() (Data, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, err
	}
	var value []Data
	for i := 0; i < size; i++ {
		data, err := r.ReadData()
		if err != nil {
			return nil, err
		}
		value = append(value, data)
	}
	return Array{Val: value}, nil
}
