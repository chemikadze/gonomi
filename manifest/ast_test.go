package manifest

import (
	"testing"
)

type typeCase struct {
	DataType DataType
	Result   string
}

func TestComponentTypes(t *testing.T) {
	cases := []typeCase{
		typeCase{StringDataType{}, "string"},
		typeCase{IntDataType{}, "int"},
		typeCase{BoolDataType{}, "bool"},
		typeCase{ListDataType{BoolDataType{}}, "list<bool>"},
		typeCase{ListDataType{ListDataType{BoolDataType{}}}, "list<list<bool>>"},
		typeCase{MapDataType{IntDataType{}, BoolDataType{}}, "map<int, bool>"},
		typeCase{RecordDataType{map[string]DataType{
			"a": IntDataType{},
			"b": ListDataType{BoolDataType{}}}}, "record<int a, list<bool> b>"},
	}
	for _, el := range cases {
		if el.DataType.DataTypeName() != el.Result {
			t.Error(el.DataType.DataTypeName(), el.Result)
		}
	}
}

func TestSimpleDeclaration(t *testing.T) {
	app := Application{
		CompositeComponent{
			Type{CompositeTypeName},
			Configuration{},
			map[string]Component{
				"main": LeafComponent{
					Type{"test.Component"},
					Configuration{
						"test.scalar":    "azaza",
						"test.composite": []string{"azaza"},
					},
					map[string]LeafInterface{
						"myinterface": LeafInterface{
							"mypin": DirectedPinType{
								Sends,
								SignalPin{StringDataType{}},
							},
						},
					}},
				"other": LeafComponent{
					Type{"test.Component"},
					Configuration{},
					map[string]LeafInterface{}},
			},
			map[string]CompositeInterface{
				"myinterface": CompositeInterface{
					"mypin": PinBinding{
						"main",
						PinId{"myinterface", "mypin"},
					},
				},
			},
			[]Binding{
				Binding{
					ComponentBindingTarget{"other"},
					InterfaceBindingTarget{"main", "myinterface"}}}}}
	if &app == nil {
		t.Error("WAT") // mock condition to work around unused warning
	}
}
