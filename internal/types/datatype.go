package types

import "strings"

type DataType uint8

//NOTE: for now let just use nil pointer for NULL and see how far we get

const (
	TYPE_INT DataType = iota + 1
	TYPE_FLOAT
	TYPE_BOOL
	TYPE_TEXT
)

func ParseDataType(typeStr string) DataType {
	switch strings.ToUpper(typeStr) {
	case "INT":
		return TYPE_INT
	case "FLOAT":
		return TYPE_FLOAT
	case "BOOL":
		return TYPE_BOOL
	case "TEXT":
		return TYPE_TEXT
	default:
		return 0
	}
}

func (dt DataType) String() string {
	switch dt {
	case TYPE_INT:
		return "INT"
	case TYPE_FLOAT:
		return "FLOAT"
	case TYPE_BOOL:
		return "BOOL"
	case TYPE_TEXT:
		return "TEXT"
	default:
		return "UNKNOWN"
	}
}
