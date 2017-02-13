package resp

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", func() {
	Describe("#WriteSimpleString", func() {
		It("should write a simple string", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteSimpleString("a simple string")
			Expect(err).To(BeNil())
			p := make([]byte, 18)
			b.Read(p)
			Expect(string(p)).To(Equal("+a simple string\r\n"))
		})
	})
	Describe("#WriteError", func() {
		It("should write an error", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteError("an error")
			Expect(err).To(BeNil())
			p := make([]byte, 11)
			b.Read(p)
			Expect(string(p)).To(Equal("-an error\r\n"))
		})
	})
	Describe("#WriteBulkString", func() {
		It("should write a bulk string", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			bs := "a bulk string 32 chars in length"
			err := w.WriteBulkString(&bs)
			Expect(err).To(BeNil())
			p := make([]byte, 39)
			b.Read(p)
			Expect(string(p)).To(Equal("$32\r\na bulk string 32 chars in length\r\n"))
		})
		It("should write a null bulk string", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteBulkString(nil)
			Expect(err).To(BeNil())
			p := make([]byte, 5)
			b.Read(p)
			Expect(string(p)).To(Equal("$-1\r\n"))
		})
	})
	Describe("#WriteInteger", func() {
		It("should write an integer", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteInteger(12345)
			Expect(err).To(BeNil())
			p := make([]byte, 8)
			b.Read(p)
			Expect(string(p)).To(Equal(":12345\r\n"))
		})
		It("should write a negative integer", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteInteger(-12345)
			Expect(err).To(BeNil())
			p := make([]byte, 9)
			b.Read(p)
			Expect(string(p)).To(Equal(":-12345\r\n"))
		})
	})
	Describe("#WriteArray", func() {
		It("should write an array", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			bs := "bulk string"
			err := w.WriteArray(SimpleString{"OK"}, Error{"ERR invalid"}, BulkString{&bs}, BulkString{nil}, Integer{12345})
			Expect(err).To(BeNil())
			p := make([]byte, 54)
			b.Read(p)
			Expect(string(p)).To(Equal("*5\r\n+OK\r\n-ERR invalid\r\n$11\r\nbulk string\r\n$-1\r\n:12345\r\n"))
		})
	})
	Describe("#WriteData", func() {
		It("should write a simple string", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteData(SimpleString{"a simple string"})
			Expect(err).To(BeNil())
			p := make([]byte, 18)
			b.Read(p)
			Expect(string(p)).To(Equal("+a simple string\r\n"))
		})
		It("should write an error", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteData(Error{"an error"})
			Expect(err).To(BeNil())
			p := make([]byte, 11)
			b.Read(p)
			Expect(string(p)).To(Equal("-an error\r\n"))
		})
		It("should write a bulk string", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			bs := "a bulk string 32 chars in length"
			err := w.WriteData(BulkString{&bs})
			Expect(err).To(BeNil())
			p := make([]byte, 39)
			b.Read(p)
			Expect(string(p)).To(Equal("$32\r\na bulk string 32 chars in length\r\n"))
		})
		It("should write an integer", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			err := w.WriteData(Integer{12345})
			Expect(err).To(BeNil())
			p := make([]byte, 8)
			b.Read(p)
			Expect(string(p)).To(Equal(":12345\r\n"))
		})
		It("should write an array", func() {
			b := bytes.Buffer{}
			w := NewWriter(&b)
			bs := "bulk string"
			err := w.WriteData(Array{[]Data{SimpleString{"OK"}, Error{"ERR invalid"}, BulkString{&bs}, BulkString{nil}, Integer{12345}}})
			Expect(err).To(BeNil())
			p := make([]byte, 54)
			b.Read(p)
			Expect(string(p)).To(Equal("*5\r\n+OK\r\n-ERR invalid\r\n$11\r\nbulk string\r\n$-1\r\n:12345\r\n"))
		})
	})
})
