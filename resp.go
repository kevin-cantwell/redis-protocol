package resp

var (
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
)

const (
	SimpleStringPrefix byte = '+'
	ErrorsPrefix            = '-'
	IntegerPrefix           = ':'
	BulkStringPrefix        = '$'
	ArrayPrefix             = '*'
	CR                      = '\r'
	LF                      = '\n'

	CRLF = "\r\n"
)
