package manifest

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWrongYaml(t *testing.T) {
	app, err := Parse(":")
	if err == nil {
		t.Error("Error expected")
	}
	if (!app.Equal(Application{})) {
		t.Error("App should be empty application")
	}
	fmt.Println(err)
}

func TestEmptyComponents(t *testing.T) {
	app, err := Parse(`
        application:
            components: {}
    `)
	if err != nil {
		t.Error(err)
	}
	if (!app.Equal(Application{CompositeComponent{Components: make(map[string]Component)}})) {
		t.Error("App should be empty application")
	}
	fmt.Println(err)
}

func testManifest(t *testing.T, manifest string, repr Application) {
	app, err := Parse(manifest)
	if err != nil {
		t.Error(err)
	}
	expected := repr
	if !app.Equal(expected) {
		t.Errorf("\nParsed: %v\nExpect: %v", app, expected)
	}
	fmt.Println(err)
}

func TestLeafComponent(t *testing.T) {
	testManifest(t, `
        application:
            components:
                x:
                    type: test.Component
    `, Application{CompositeComponent{
		Components: map[string]Component{
			"x": LeafComponent{
				Type:          Type{"test.Component"},
				Configuration: Configuration{},
				Interfaces:    map[string]LeafInterface{},
			},
		}}})
}

func TestConfiguration(t *testing.T) {
	testManifest(t, `
        application:
            components:
                x:
                    type: test.Component
                    configuration:
                        sample.string: c
                        sample.list: [1]
                        sample.map: {3: 4}
    `, Application{CompositeComponent{
		Components: map[string]Component{
			"x": LeafComponent{
				Type: Type{"test.Component"},
				Configuration: Configuration{
					"sample.string": "c",
					"sample.list":   []interface{}{1},
					"sample.map":    map[interface{}]interface{}{3: 4},
				},
				Interfaces: map[string]LeafInterface{},
			},
		}}})
}

func TestInterface(t *testing.T) {
	testManifest(t, `
        application:
            components:
                x:
                    type: test.Component
                    interfaces:
                        myinterface:
                            mypin1: publish-signal(string)
                            mypin2: consume-signal(string)
                            mypin3: send-command()
                            mypin4: receive-command()
                        myrequired:
                            mypin: publish-signal(string)
                    required: [myrequired]
    `, Application{CompositeComponent{
		Components: map[string]Component{
			"x": LeafComponent{
				Type:          Type{"test.Component"},
				Configuration: Configuration{},
				Interfaces: map[string]LeafInterface{
					"myinterface": LeafInterface{
						Pins: map[string]DirectedPinType{
							"mypin1": {Sends, SignalPin{StringDataType{}}},
							"mypin2": {Receives, SignalPin{StringDataType{}}},
							"mypin3": {Sends, CommandPin{RecordDataType{}, RecordDataType{}, RecordDataType{}}},
							"mypin4": {Receives, CommandPin{}},
						},
					},
					"myrequired": LeafInterface{
						Pins: map[string]DirectedPinType{
							"mypin": {Sends, SignalPin{StringDataType{}}},
						},
						Required: true,
					},
				},
			},
		}}})
}

func TestPinTypes(t *testing.T) {
	cases := map[string]DataType{
		"int":                IntDataType{},
		"string":             StringDataType{},
		"bool":               BoolDataType{},
		"list<string>":       ListDataType{StringDataType{}},
		"list<list<string>>": ListDataType{ListDataType{StringDataType{}}},
		// "list< list <string> >": ListDataType{ListDataType{StringDataType{}}},
		// "map<string, string>":   MapDataType{StringDataType{}, StringDataType{}},
	}
	for manifest, expected := range cases {
		parsed, err := parseDataType(manifest)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(parsed, expected) {
			t.Errorf("\nParsed: %v\nExpect: %v", parsed, expected)
		}
	}
}
