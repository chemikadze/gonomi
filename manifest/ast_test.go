package manifest

import (
	"github.com/chemikadze/gonomi/manifest/datatype"
	"testing"
)

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
							Pins: map[string]DirectedPinType{
								"mypin": {Sends, SignalPin{datatype.String{}}},
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
