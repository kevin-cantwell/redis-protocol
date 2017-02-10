package resp

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	Describe("#ReadData", func() {
		It("should read simple strings", func() {
			r := NewReader(strings.NewReader("+This is a simple string\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			data, ok := d.(SimpleString)
			Expect(ok).To(BeTrue())
			Expect(data.Val).To(Equal("This is a simple string"))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read errors", func() {
			r := NewReader(strings.NewReader("-This is an error\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			data, ok := d.(Error)
			Expect(ok).To(BeTrue())
			Expect(data.Val).To(Equal("This is an error"))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read integers", func() {
			r := NewReader(strings.NewReader(":1234567\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			data, ok := d.(Integer)
			Expect(ok).To(BeTrue())
			Expect(data.Val).To(Equal(int64(1234567)))

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
			data, ok := d.(BulkString)
			Expect(ok).To(BeTrue())
			Expect(*data.Val).To(Equal("foo bar\r\nbaz"))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read NULL bulk strings", func() {
			r := NewReader(strings.NewReader("$-1\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			data, ok := d.(BulkString)
			Expect(ok).To(BeTrue())
			Expect(data.Val).To(BeNil())

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read arrays", func() {
			r := NewReader(strings.NewReader("*5\r\n+OK\r\n-ERR foo\r\n:2\r\n$3\r\nbar\r\n$-1\r\n"))
			d, err := r.ReadData()
			Expect(err).To(BeNil())
			a, ok := d.(Array)
			if !ok {
				Fail("Expected type to be Array")
			}
			value := a.Val
			Expect(len(value)).To(Equal(5))
			value0, ok := value[0].(SimpleString)
			if !ok {
				Fail("Expected type to be SimpleString")
			}
			Expect(value0.Val).To(Equal("OK"))
			value1, ok := value[1].(Error)
			if !ok {
				Fail("Expected type to be Error")
			}
			Expect(value1.Val).To(Equal("ERR foo"))
			value2, ok := value[2].(Integer)
			if !ok {
				Fail("Expected type to be Integer")
			}
			Expect(value2.Val).To(Equal(int64(2)))
			value3, ok := value[3].(BulkString)
			if !ok {
				Fail("Expected type to be BulkString")
			}
			Expect(*value3.Val).To(Equal("bar"))
			value4, ok := value[4].(BulkString)
			if !ok {
				Fail("Expected type to be BulkString")
			}
			Expect(value4.Val).To(BeNil())

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})
	})
})
