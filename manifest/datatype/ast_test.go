package datatype

import (
	"testing"
)

type typeCase struct {
	DataType DataType
	Results  []string
}

func TestDataTypes(t *testing.T) {
	cases := []typeCase{
		typeCase{String{}, []string{"string"}},
		typeCase{Int{}, []string{"int"}},
		typeCase{Bool{}, []string{"bool"}},
		typeCase{List{Bool{}}, []string{"list<bool>"}},
		typeCase{List{List{Bool{}}}, []string{"list<list<bool>>"}},
		typeCase{Map{Int{}, Bool{}}, []string{"map<int, bool>"}},
		typeCase{
			Record{map[string]DataType{
				"a": Int{},
				"b": List{Bool{}}},
			},
			[]string{"record<int a, list<bool> b>", "record<list<bool> b, int a>"}},
	}
	for _, el := range cases {
		name := el.DataType.DataTypeName()
		ok := false
		for _, element := range el.Results {
			if name == element {
				ok = true
				break
			}
		}
		if !ok {
			t.Error(name, el.Results[0])
		}
	}
}
