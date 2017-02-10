package resp

const (
	simpleStringPrefix byte = '+'
	errorPrefix             = '-'
	integerPrefix           = ':'
	bulkStringPrefix        = '$'
	arrayPrefix             = '*'
)

type Data interface {
	protected() // Only this package may implement
}

type Error struct {
	Val string
}

func (d Error) protected() {}

type SimpleString struct {
	Val string
}

func (d SimpleString) protected() {}

type BulkString struct {
	Val *string // Bulk strings support nil values
}

func (d BulkString) protected() {}

type Integer struct {
	Val int64
}

func (d Integer) protected() {}

type Array struct {
	Val []Data
}

func (d Array) protected() {}
