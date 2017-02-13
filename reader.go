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
		prefix, err := r.ReadPrefix()
		if err != nil {
			return nil, err
		}
		switch prefix {
		case '+':
			return r.ReadSimpleString()
		case '-':
			return r.ReadError()
		case ':':
			return r.ReadInteger()
		case '$':
			return r.ReadBulkString()
		case '*':
			return r.ReadArray()
		}
	}
}

func (r *Reader) ReadPrefix() (byte, error) {
	prefix, err := r.src.ReadByte()
	if err != nil {
		return 0, err
	}
	switch prefix {
	case '+', '-', ':', '$', '*':
		return prefix, nil
	default:
		return 0, fmt.Errorf("resp: unknown prefix %q", prefix)
	}
}

func (r *Reader) ReadSimpleString() (SimpleString, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return "", err
	}
	return SimpleString(line), nil
}

func (r *Reader) ReadError() (Error, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return "", err
	}
	return Error(line), nil
}

func (r *Reader) ReadInteger() (Integer, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}
	return Integer(i), nil
}

func (r *Reader) ReadBulkString() (BulkString, error) {
	line1, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(string(line1))
	if err != nil {
		return nil, err
	}
	if size == -1 {
		return nil, nil
	}
	line2 := make([]byte, size)
	for i := 0; i < size; i++ { // The extra 2 bytes is for the line terminator
		b, err := r.src.ReadByte()
		if err != nil {
			return nil, err
		}
		line2[i] = b
	}
	// Make sure to read the terminating CRLF
	if _, _, err := r.src.ReadLine(); err != nil {
		return nil, err
	}
	return BulkString(line2), nil
}

func (r *Reader) ReadArray() (Array, error) {
	line, _, err := r.src.ReadLine()
	if err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, err
	}
	var value Array
	for i := 0; i < size; i++ {
		data, err := r.ReadData()
		if err != nil {
			return nil, err
		}
		value = append(value, data)
	}
	return value, nil
}
