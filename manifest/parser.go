package manifest

import (
	"fmt"
	"github.com/chemikadze/gonomi/manifest/datatype"
	"github.com/chemikadze/gonomi/manifest/parsing"
	"gopkg.in/yaml.v2"
	"strings"
)

var _ = fmt.Printf

type applicationRoot struct {
	Application application `yaml:application`
}

type application struct {
	Components map[string]component `yaml:components`
	Bindings   [][]string           `yaml:bindings`
}

type component struct {
	Type          string                                      `yaml:type`
	Configuration map[interface{}]interface{}                 `yaml:configuration`
	Interfaces    map[interface{}]map[interface{}]interface{} `yaml:interfaces`
	Required      []interface{}                               `yaml:required`
	Bindings      []interface{}                               `yaml:bindings`
}

func Parse(manifest string) (Application, error) {
	// TODO: composite components
	// TODO: bindings for nested components
	m := applicationRoot{}
	err := yaml.Unmarshal([]byte(manifest), &m)
	if err != nil {
		return Application{}, err
	}
	components, err := parseComponents(m.Application.Components)
	if err != nil {
		return Application{}, err
	}
	bindings, err := parseBindings(m.Application.Bindings)
	return Application{CompositeComponent{Components: components, Bindings: bindings}}, err
}

func parseComponents(components map[string]component) (map[string]Component, error) {
	if len(components) == 0 {
		return nil, nil
	}
	acc := make(map[string]Component)
	for id, component := range components {
		interfaces, err := yamlMapToInterfacesMap(component.Interfaces)
		if err != nil {
			return nil, err
		}
		leafComponent := LeafComponent{
			Type:          Type{component.Type},
			Configuration: Configuration(yamlMapToSimpleMap(component.Configuration)),
			Interfaces:    interfaces,
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

func yamlMapToInterfacesMap(original map[interface{}]map[interface{}]interface{}) (map[string]LeafInterface, error) {
	result := make(map[string]LeafInterface)
	for k, v := range original {
		switch t := k.(type) {
		case string:
			value, err := parseLeafInterface(v)
			if err != nil {
				return nil, err
			}
			result[t] = value
		}
	}
	return result, nil
}

func parseLeafInterface(original map[interface{}]interface{}) (LeafInterface, error) {
	result := make(map[string]DirectedPinType)
	for k, v := range original {
		switch k := k.(type) {
		case string:
			switch v := v.(type) {
			case string:
				pinType, err := parseDirectedPinType(v)
				if err != nil {
					return LeafInterface{}, err
				}
				result[k] = pinType
			}
		}
	}
	return LeafInterface{result, false}, nil
}

func parseDirectedPinType(repr string) (DirectedPinType, error) {
	pinAndTypes := strings.SplitN(repr, "(", 2)
	pinAndTypes[1] = pinAndTypes[1][:len(pinAndTypes[1])-1]
	switch pinAndTypes[0] {
	case "publish-signal":
		t, _ := datatype.Parse(pinAndTypes[1])
		return DirectedPinType{Sends, SignalPin{t}}, nil
	case "consume-signal":
		t, _ := datatype.Parse(pinAndTypes[1])
		return DirectedPinType{Receives, SignalPin{t}}, nil
	case "send-command": // FIXME: correct parameter parsing
		return DirectedPinType{Sends, CommandPin{datatype.Record{}, datatype.Record{}, datatype.Record{}}}, nil
	case "receive-command": // FIXME: correct parameter parsing
		return DirectedPinType{Receives, CommandPin{datatype.Record{}, datatype.Record{}, datatype.Record{}}}, nil
	default:
		return DirectedPinType{}, parsing.ManifestError{fmt.Sprintf("Unknown pin type: %s", pinAndTypes[0]), 0, 0}
	}
}

func parseBindings(repr [][]string) ([]Binding, error) {
	if len(repr) == 0 {
		return nil, nil
	}
	bindings := make([]Binding, 0, len(repr))
	for _, binding := range repr {
		if len(binding) != 2 {
			return nil, parsing.ManifestError{fmt.Sprintf("Expected list of two for binding, got: %s", binding), 0, 0}
		}
		binding, err := parseBinding(binding[0], binding[1])
		if err != nil {
			return nil, err
		}
		bindings = append(bindings, binding)
	}
	return bindings, nil
}

func parseBinding(left, right string) (Binding, error) {
	first, err := parseBindingTarget(left)
	if err != nil {
		return Binding{}, err
	}
	second, err := parseBindingTarget(right)
	if err != nil {
		return Binding{}, err
	}
	return Binding{first, second}, nil
}

func parseBindingTarget(target string) (BindingTarget, error) {
	splitted := strings.Split(target, "#")
	if len(splitted) == 1 {
		return ComponentBindingTarget{ComponentId{strings.Split(splitted[0], ".")}}, nil
	} else if len(splitted) == 2 {
		return InterfaceBindingTarget{ComponentId{strings.Split(splitted[0], ".")}, splitted[1]}, nil
	} else {
		return nil, parsing.ManifestError{fmt.Sprintf("Unexpected binding target: %s", target), 0, 0}
	}
}
