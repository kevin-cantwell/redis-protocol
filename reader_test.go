package resp

import (
	"fmt"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	Describe("#ReadPrefix", func() {
		It("should read valid prefixes", func() {
			r := NewReader(strings.NewReader("*4\r\n+OK\r\n-ERR foo\r\n:123\r\n$3\r\nbar\r\n"))
			for _, expected := range []byte{'*', '+', '-', ':', '$'} {
				Expect(r.ReadPrefix()).To(Equal(expected))
				r.src.ReadLine()
			}
		})
		It("should reject invalid prefixes", func() {
			r := NewReader(strings.NewReader("#nope\r\n=nope\r\n&nope\r\n%nope\r\n@nope\r\n"))
			for _, expected := range []byte{'#', '=', '&', '%', '@'} {
				_, err := r.ReadPrefix()
				Expect(err.Error()).To(Equal(fmt.Sprintf("resp: unknown prefix %q", expected)))
				r.src.ReadLine()
			}
		})
	})
	Describe("#ReadSimpleString", func() {
		It("should read a simple string", func() {
			r := NewReader(strings.NewReader("+OK\r\n"))
			r.ReadPrefix()
			Expect(r.ReadSimpleString()).To(Equal(SimpleString("OK")))
		})
	})
	Describe("#ReadError", func() {
		It("should read an error", func() {
			r := NewReader(strings.NewReader("-ERR foo\r\n"))
			r.ReadPrefix()
			Expect(r.ReadError()).To(Equal(Error("ERR foo")))
		})
	})
	Describe("#ReadBulkString", func() {
		It("should read a bulk string", func() {
			r := NewReader(strings.NewReader("$12\r\nfoo bar\r\nbaz\r\n"))
			r.ReadPrefix()
			Expect(r.ReadBulkString()).To(Equal(BulkString("foo bar\r\nbaz")))
		})
		It("should read a null bulk string", func() {
			r := NewReader(strings.NewReader("$-1\r\n"))
			r.ReadPrefix()
			Expect(r.ReadBulkString()).To(Equal(BulkString(nil)))
		})
		It("should know the difference between nulls and empty strings", func() {
			Expect(BulkString("")).ToNot(Equal(BulkString(nil)))
		})
	})
	Describe("#ReadInteger", func() {
		It("should read an integer", func() {
			r := NewReader(strings.NewReader(":-12345\r\n"))
			r.ReadPrefix()
			Expect(r.ReadInteger()).To(Equal(Integer(-12345)))
		})
	})
	Describe("#ReadArray", func() {
		It("should read an array", func() {
			r := NewReader(strings.NewReader("*5\r\n+OK\r\n-ERR foo\r\n:123\r\n$3\r\nbar\r\n$-1\r\n"))
			r.ReadPrefix()
			Expect(r.ReadArray()).To(Equal(Array{SimpleString("OK"), Error("ERR foo"), Integer(123), BulkString("bar"), BulkString(nil)}))
		})
	})
	Describe("#ReadData", func() {
		It("should read simple strings", func() {
			r := NewReader(strings.NewReader("+This is a simple string\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			Expect(d).To(Equal(SimpleString("This is a simple string")))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read errors", func() {
			r := NewReader(strings.NewReader("-This is an error\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			Expect(d).To(Equal(Error("This is an error")))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read integers", func() {
			r := NewReader(strings.NewReader(":1234567\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			Expect(d).To(Equal(Integer(1234567)))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should validate integers", func() {
			r := NewReader(strings.NewReader(":not an integer\r\n"))
			_, err := r.ReadData()
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(`strconv.ParseInt: parsing "not an integer": invalid syntax`))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read bulk strings", func() {
			r := NewReader(strings.NewReader("$12\r\nfoo bar\r\nbaz\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			Expect(d).To(Equal(BulkString("foo bar\r\nbaz")))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read NULL bulk strings", func() {
			r := NewReader(strings.NewReader("$-1\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			Expect(d).To(BeNil())

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read arrays", func() {
			r := NewReader(strings.NewReader("*5\r\n+OK\r\n-ERR foo\r\n:2\r\n$3\r\nbar\r\n$-1\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			value, ok := d.(Array)
			if !ok {
				Fail("Expected type to be Array")
			}
			Expect(len(value)).To(Equal(5))
			Expect(value[0]).To(Equal(SimpleString("OK")))
			Expect(value[1]).To(Equal(Error("ERR foo")))
			Expect(value[2]).To(Equal(Integer(2)))
			Expect(value[3]).To(Equal(BulkString("bar")))
			Expect(value[4]).To(BeNil())

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})
	})
})
