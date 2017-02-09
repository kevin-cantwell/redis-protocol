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
			data, err := r.ReadData()
			Expect(err).To(BeNil())
			_, ok := data.(SimpleString)
			Expect(ok).To(BeTrue())
			Expect(data.Value()).To(Equal("This is a simple string"))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read errors", func() {
			r := NewReader(strings.NewReader("-This is an error\r\n"))
			data, err := r.ReadData()
			Expect(err).To(BeNil())
			_, ok := data.(Error)
			Expect(ok).To(BeTrue())
			Expect(data.Value()).To(Equal("This is an error"))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read integers", func() {
			r := NewReader(strings.NewReader(":1234567\r\n"))
			data, err := r.ReadData()
			Expect(err).To(BeNil())
			_, ok := data.(Integer)
			Expect(ok).To(BeTrue())
			Expect(data.Value()).To(Equal(int64(1234567)))

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
			data, err := r.ReadData()
			Expect(err).To(BeNil())
			_, ok := data.(BulkString)
			Expect(ok).To(BeTrue())
			Expect(data.Value()).To(Equal("foo bar\r\nbaz"))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})

		It("should read arrays", func() {
			r := NewReader(strings.NewReader("*5\r\n+OK\r\n-ERR foo\r\n:2\r\n$3\r\nbar\r\n$-1\r\n"))
			data, err := r.ReadData()
			Expect(err).To(BeNil())
			_, ok := data.(Array)
			if _, ok := data.(Array); !ok {
				Fail("Expected type to be Array")
			}
			value, ok := data.Value().([]Data)
			if !ok {
				Fail("Expected type to be []Data")
			}
			Expect(len(value)).To(Equal(5))
			if _, ok := value[0].(SimpleString); !ok {
				Fail("Expected type to be SimpleString")
			}
			Expect(value[0].Value()).To(Equal("OK"))
			if _, ok := value[1].(Error); !ok {
				Fail("Expected type to be Error")
			}
			Expect(value[1].Value()).To(Equal("ERR foo"))
			if _, ok := value[2].(Integer); !ok {
				Fail("Expected type to be Integer")
			}
			Expect(value[2].Value()).To(Equal(int64(2)))
			if _, ok := value[3].(BulkString); !ok {
				Fail("Expected type to be BulkString")
			}
			Expect(value[3].Value()).To(Equal("bar"))
			if _, ok := value[4].(BulkString); !ok {
				Fail("Expected type to be BulkString")
			}
			Expect(value[4].Value()).To(BeNil())

			// Expect(strconv.Quote(data.RawString())).To(Equal(`"OK\nERR foo\n2\nbar\n\n"`))
			// Expect(strconv.Quote(data.String())).To(Equal(`"OK\n(error) ERR foo\n2\nbar\n(nil)\n"`))

			// No more bytes should be left in buffer
			_, err = r.src.ReadByte()
			Expect(err).To(Equal(io.EOF))
		})
	})
})
