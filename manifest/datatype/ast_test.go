package datatype

import (
	"testing"
)

type typeCase struct {
	DataType DataType
	Result   string
}

func TestDataTypes(t *testing.T) {
	cases := []typeCase{
		typeCase{String{}, "string"},
		typeCase{Int{}, "int"},
		typeCase{Bool{}, "bool"},
		typeCase{List{Bool{}}, "list<bool>"},
		typeCase{List{List{Bool{}}}, "list<list<bool>>"},
		typeCase{Map{Int{}, Bool{}}, "map<int, bool>"},
		typeCase{Record{map[string]DataType{
			"a": Int{},
			"b": List{Bool{}}}}, "record<int a, list<bool> b>"},
	}
	for _, el := range cases {
		if el.DataType.DataTypeName() != el.Result {
			t.Error(el.DataType.DataTypeName(), el.Result)
		}
	}
}
