package datatype

import (
	"strings"
)

// data types
type DataType interface {
	DataTypeName() string
}

type String struct{}

func (String) DataTypeName() string {
	return "string"
}

type Int struct{}

func (Int) DataTypeName() string {
	return "int"
}

type Bool struct{}

func (Bool) DataTypeName() string {
	return "bool"
}

type List struct {
	ElementDataType DataType
}

func (l List) DataTypeName() string {
	return "list<" + l.ElementDataType.DataTypeName() + ">"
}

type Map struct {
	KeyDataType   DataType
	ValueDataType DataType
}

func (m Map) DataTypeName() string {
	return "map<" + m.KeyDataType.DataTypeName() + ", " + m.ValueDataType.DataTypeName() + ">"
}

type Record struct {
	Fields map[string]DataType
}

func (r Record) DataTypeName() string {
	items := make([]string, 0, 5)
	for key, value := range r.Fields {
		items = append(items, value.DataTypeName()+" "+key)
	}
	return "record<" + strings.Join(items, ", ") + ">"
}
