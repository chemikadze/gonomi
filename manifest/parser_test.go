package manifest

import (
	"errors"
	"fmt"
	"github.com/chemikadze/gonomi/manifest/datatype"
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
	if (!app.Equal(Application{})) {
		t.Error("Empty manifest parsing should yield empty application, got:", app)
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

func testManifestError(t *testing.T, manifest string, expectedErr error) {
	app, err := Parse(manifest)
	if err == nil {
		t.Errorf("Expected error %s, got result %s", expectedErr, app)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("\nRaised: %v\nExpect: %v", err.Error(), expectedErr.Error())
	}
	if (!app.Equal(Application{})) {
		t.Error("Failed parsing should yield empty application, got:", app)
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
							"mypin1": {Sends, SignalPin{datatype.String{}}},
							"mypin2": {Receives, SignalPin{datatype.String{}}},
							"mypin3": {Sends, CommandPin{datatype.Record{}, datatype.Record{}, datatype.Record{}}},
							"mypin4": {Receives, CommandPin{}},
						},
					},
					"myrequired": LeafInterface{
						Pins: map[string]DirectedPinType{
							"mypin": {Sends, SignalPin{datatype.String{}}},
						},
						Required: true,
					},
				},
			},
		}}})
}

func TestInterfaceError(t *testing.T) {
	testManifestError(t, `
        application:
            components:
                x:
                    type: test.Component
                    interfaces:
                        myinterface:
                            mypin1: plubish-signal(string)
    `,
		errors.New("Unknown pin type: plubish-signal"))
}

func TestBindings(t *testing.T) {
	testManifest(t, `
        application:
            components:
                x:
                    type: test.Component
                y:
                    type: test.Component
            bindings:
                - [x, y]
                - [x#i, y]
    `, Application{CompositeComponent{
		Components: map[string]Component{
			"x": LeafComponent{
				Type:          Type{"test.Component"},
				Configuration: Configuration{},
				Interfaces:    map[string]LeafInterface{},
			},
			"y": LeafComponent{
				Type:          Type{"test.Component"},
				Configuration: Configuration{},
				Interfaces:    map[string]LeafInterface{},
			},
		},
		Bindings: []Binding{
			{ComponentBindingTarget{ComponentId{[]string{"x"}}}, ComponentBindingTarget{ComponentId{[]string{"y"}}}},
			{InterfaceBindingTarget{ComponentId{[]string{"x"}}, "i"}, ComponentBindingTarget{ComponentId{[]string{"y"}}}},
		}}})
}
