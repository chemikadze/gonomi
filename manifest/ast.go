package manifest

import (
	"github.com/chemikadze/gonomi/manifest/datatype"
	"reflect"
)

const (
	CompositeTypeName = "core.Composite"
)

type Type struct {
	Name string
}

type PinId struct {
	Interface string
	Pin       string
}

type Component interface {
	GetType() Type
	GetConfiguration() map[string]interface{}
}

type Configuration map[string]interface{}

// component with children
type CompositeComponent struct {
	Type          Type
	Configuration Configuration
	Components    map[string]Component
	Interfaces    map[string]CompositeInterface
	Bindings      []Binding
}

func (c CompositeComponent) GetType() Type {
	return c.Type
}

func (c CompositeComponent) GetConfiguration() map[string]interface{} {
	return c.Configuration
}

type CompositeInterface map[string]PinBinding

type PinBinding struct {
	TargetComponent string
	TargetPin       PinId
}

type Binding struct {
	Left  BindingTarget
	Right BindingTarget
}

// side of binding
type BindingTarget interface {
	// hack
	IsBindingTarget()
}

type ComponentBindingTarget struct {
	Component string
}

func (c ComponentBindingTarget) IsBindingTarget() {}

type InterfaceBindingTarget struct {
	Component string
	Interface string
}

func (c InterfaceBindingTarget) IsBindingTarget() {}

type Application struct {
	CompositeComponent
}

func (c Application) Equal(other Application) bool {
	return reflect.DeepEqual(c, other)
}

// component without children
type LeafComponent struct {
	Type          Type
	Configuration Configuration
	Interfaces    map[string]LeafInterface
}

func (c LeafComponent) GetType() Type {
	return c.Type
}

func (c LeafComponent) GetConfiguration() map[string]interface{} {
	return c.Configuration
}

type LeafInterface struct {
	Pins     map[string]DirectedPinType
	Required bool
}

// pin in interface
type DirectedPinType struct {
	Direction Direction
	PinType   PinType
}

type Direction interface {
	IsSend() bool
	IsReceive() bool
}

type direction bool

var (
	Sends    Direction = direction(true)
	Receives Direction = direction(false)
)

func (d direction) IsSend() bool {
	return bool(d)
}

func (d direction) IsReceive() bool {
	return bool(!d)
}

type PinType interface {
	PinTypeName() string
}

type SignalPin struct {
	DataType datatype.DataType
}

func (s SignalPin) PinTypeName() string {
	return "signal"
}

type ConfigurationPin struct {
	DataType datatype.DataType
}

func (s ConfigurationPin) PinTypeName() string {
	return "configuration"
}

type CommandPin struct {
	Arguments datatype.Record
	Progress  datatype.Record
	Result    datatype.Record
}

func (s CommandPin) PinTypeName() string {
	return "command"
}
