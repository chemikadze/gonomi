package manifest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

var _ = fmt.Printf

type ManifestError struct {
	message string
	line    int
	column  int
}

func (e ManifestError) Error() string {
	return e.message
}

type applicationRoot struct {
	Application application `yaml:application`
}

type application struct {
	Components map[string]component `yaml:components`
}

type component struct {
	Type          string                                      `yaml:type`
	Configuration map[interface{}]interface{}                 `yaml:configuration`
	Interfaces    map[interface{}]map[interface{}]interface{} `yaml:interfaces`
	Required      []interface{}                               `yaml:required`
}

func Parse(manifest string) (Application, error) {
	m := applicationRoot{}
	err := yaml.Unmarshal([]byte(manifest), &m)
	if err != nil {
		return Application{}, err
	}
	components, err := parseComponents(m.Application.Components)
	if err != nil {
		return Application{}, err
	}
	return Application{CompositeComponent{Components: components}}, err
}

func parseComponents(components map[string]component) (map[string]Component, error) {
	acc := make(map[string]Component)
	for id, component := range components {
		leafComponent := LeafComponent{
			Type:          Type{component.Type},
			Configuration: Configuration(yamlMapToSimpleMap(component.Configuration)),
			Interfaces:    yamlMapToInterfacesMap(component.Interfaces),
		}
		for _, reqname := range component.Required {
			switch reqname := reqname.(type) {
			case string:
				iface := leafComponent.Interfaces[reqname] // oh my, by-value gets!
				iface.Required = true
				leafComponent.Interfaces[reqname] = iface
			}
		}
		acc[id] = leafComponent
	}
	return acc, nil
}

func yamlMapToSimpleMap(original map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range original {
		switch t := k.(type) {
		case string:
			result[t] = v
		}
	}
	return result
}

func yamlMapToInterfacesMap(original map[interface{}]map[interface{}]interface{}) map[string]LeafInterface {
	result := make(map[string]LeafInterface)
	for k, v := range original {
		switch t := k.(type) {
		case string:
			result[t] = parseLeafInterface(v)
		}
	}
	return result
}

func parseLeafInterface(original map[interface{}]interface{}) LeafInterface {
	result := make(map[string]DirectedPinType)
	for k, v := range original {
		switch k := k.(type) {
		case string:
			switch v := v.(type) {
			case string:
				pinType, _ := parseDirectedPinType(v) // FIXME
				result[k] = pinType
			}
		}
	}
	return LeafInterface{result, false}
}

func parseDirectedPinType(repr string) (DirectedPinType, error) {
	pinAndTypes := strings.SplitN(repr, "(", 2)
	pinAndTypes[1] = pinAndTypes[1][:len(pinAndTypes[1])-1]
	switch pinAndTypes[0] {
	case "publish-signal":
		t, _ := parseDataType(pinAndTypes[1])
		return DirectedPinType{Sends, SignalPin{t}}, nil
	case "consume-signal":
		t, _ := parseDataType(pinAndTypes[1])
		return DirectedPinType{Receives, SignalPin{t}}, nil
	case "send-command": // FIXME
		return DirectedPinType{Sends, CommandPin{RecordDataType{}, RecordDataType{}, RecordDataType{}}}, nil
	case "receive-command": // FIXME
		return DirectedPinType{Receives, CommandPin{RecordDataType{}, RecordDataType{}, RecordDataType{}}}, nil
	}
	return DirectedPinType{}, ManifestError{"error", 0, 0}
}

func parseDataType(repr string) (DataType, error) {
	typenameAndTypeParameters := strings.SplitN(repr, "<", 2)
	switch typenameAndTypeParameters[0] {
	case "int":
		return IntDataType{}, nil
	case "bool":
		return BoolDataType{}, nil
	case "string":
		return StringDataType{}, nil
	case "list":
		param := typenameAndTypeParameters[1]
		typeParameter, err := parseDataType(param[:len(param)-1])
		return ListDataType{typeParameter}, err
		// FIXME map, record
	}
	return DataType(nil), ManifestError{"Unknown type " + repr, 0, 0}
}
