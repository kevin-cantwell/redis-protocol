package resp

const (
	simpleStringPrefix byte = '+'
	errorPrefix             = '-'
	integerPrefix           = ':'
	bulkStringPrefix        = '$'
	arrayPrefix             = '*'
)

type Data interface {
	Value() interface{}
	// Vals() []string
}

type Error struct {
	Val string
}

func (d Error) Prefix() byte {
	return errorPrefix
}

func (d Error) Value() interface{} {
	return d.Val
}

type SimpleString struct {
	Val string
}

func (d SimpleString) Prefix() byte {
	return simpleStringPrefix
}

func (d SimpleString) Value() interface{} {
	return d.Val
}

type BulkString struct {
	Val *string // Bulk strings support nil values
}

func (d BulkString) Prefix() byte {
	return bulkStringPrefix
}

func (d BulkString) Value() interface{} {
	if d.Val == nil {
		return nil
	}
	return *d.Val
}

type Integer struct {
	Val int64
}

func (d Integer) Prefix() byte {
	return integerPrefix
}

func (d Integer) Value() interface{} {
	return d.Val
}

type Array struct {
	Val []Data
}

func (d Array) Prefix() byte {
	return arrayPrefix
}

func (d Array) Value() interface{} {
	return d.Val
}
