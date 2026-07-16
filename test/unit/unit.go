package unit

import (
	"testing"

	"github.com/mattmunz/designlanguage/model"
	"github.com/stretchr/testify/require"
)

func newIntType(t *testing.T) model.Type {
	typ, err := model.NewType("Int", false)
	require.NoError(t, err)
	return typ
}

func NewPlane(t *testing.T) model.Component {
	p, err := model.NewComponent("Plane", "", nil)
	require.NoError(t, err)
	return p
}

func NewRectangle(t *testing.T) model.Entity {
	rectangle, err := model.NewEntity(
		"Rectangle", "", nil,
		[]model.Attribute{NewAttr(t, "Length", "", "int", false), NewAttr(t, "Width", "", "int", false)},
	)
	require.NoError(t, err)
	return rectangle
}

func NewAttr(t *testing.T, name, comment, typeName string, isArray bool) model.Attribute {
	a, err := model.NewAttribute(name, comment, typeName, isArray)
	require.NoError(t, err)
	return a
}

func NewSquare(t *testing.T) model.Entity {
	square, err := model.NewEntity(
		"Square", "", NewType(t, "Rectangle"), []model.Attribute{NewAttr(t, "Area", "", "int", false)},
	)
	require.NoError(t, err)
	return square
}

func NewType(t *testing.T, name string) model.Type {
	typ, err := model.NewType(name, false)
	require.NoError(t, err)
	return typ
}

func NewLabelledValue(t *testing.T) model.Entity {
	labelledValue, err := model.NewEntity(
		"LabelledValue", "", nil,
		[]model.Attribute{NewAttr(t, "Label", "", "String", false), NewAttr(t, "Value", "", "Int", false)},
	)
	require.NoError(t, err)
	return labelledValue
}

func NewPersonRepo(t *testing.T) model.Object {
	addMethod, err := model.NewMethod(
		"Add", "", []model.Param{model.NewParam("Person", NewType(t, "Person"))}, []model.Param{},
	)
	require.NoError(t, err)

	getMethod, err := model.NewMethod("Get", "",
		[]model.Param{model.NewParam("name", NewType(t, "String"))},
		[]model.Param{model.NewParam("person", NewType(t, "Person"))})
	require.NoError(t, err)

	personRepo, err := model.NewObject(
		"PersonRepository", "", nil, []model.Attribute{},
		[]model.Method{
			addMethod,
			getMethod,
		})
	require.NoError(t, err)
	return personRepo
}

func NewSpinner(t *testing.T) model.Object {
	method, err := model.NewMethod("Spin", "",
		[]model.Param{model.NewParam("velocity", newIntType(t)), model.NewParam("duration", newIntType(t))},
		[]model.Param{})
	require.NoError(t, err)

	o, err := model.NewObject(
		"Spinner", "", nil, []model.Attribute{NewAttr(t, "Radius", "", "int", false)},
		[]model.Method{method},
	)
	require.NoError(t, err)

	return o
}

func NewSpinningSquare(t *testing.T) model.Object {
	method, err := model.NewMethod("Spin", "",
		[]model.Param{model.NewParam("Velocity", newIntType(t)), model.NewParam("Duration", newIntType(t))},
		[]model.Param{})
	require.NoError(t, err)

	o, err := model.NewObject(
		"SpinningSquare", "", NewType(t, "Square"), []model.Attribute{},
		[]model.Method{method},
	)
	require.NoError(t, err)

	return o
}

func NewColorEnum(t *testing.T) model.Enum {
	e, err := model.NewEnum("Color", "Red", "Green", "Blue")
	require.NoError(t, err)
	return e
}

func NewAppkitDesign(t *testing.T) model.Design {
	appAttrs := []model.Attribute{
		NewAttr(t, "ID", "", "String", false),
		NewAttr(t, "Version", "Semver preferred.", "String", false),
		NewAttr(t, "ConfigName", "", "String", false),
	}
	app, err := model.NewEntity("App", "A software application.", nil, appAttrs)
	require.NoError(t, err)

	exeMethod, err := model.NewMethod("Execute", "", []model.Param{}, []model.Param{})
	require.NoError(t, err)
	cmdMethods := []model.Method{
		exeMethod,
	}
	command, err := model.NewObject("Command", "", nil, []model.Attribute{}, cmdMethods)
	require.NoError(t, err)

	cliAttributes := []model.Attribute{
		NewAttr(t, "AppID", "", "String", false),
		NewAttr(t, "Name", "", "String", false),
		NewAttr(t, "ShortDescription", "", "String", false),
	}

	nrReturnVals := []model.Param{
		model.NewParam("Cmd", NewType(t, "CommandImpl")),
		model.NewParam("Err", NewType(t, "Error")),
	}
	nrMethod, err := model.NewMethod("NewRootCommand", "", []model.Param{}, nrReturnVals)
	require.NoError(t, err)

	slMethod, err := model.NewMethod("SetLogger", "", []model.Param{model.NewParam("Logger", NewType(t, "Logger"))}, []model.Param{})
	require.NoError(t, err)

	cliMethods := []model.Method{nrMethod, slMethod}
	cli, err := model.NewObject("CLI", "Command Line Interface.", NewType(t, "Command"), cliAttributes, cliMethods)
	require.NoError(t, err)

	newMethod, err := model.NewMethod("New", "", []model.Param{}, []model.Param{model.NewParam("Cmd", NewType(t, "Command"))})
	require.NoError(t, err)

	commandFactory, err := model.NewObject(
		"CommandFactory", "", nil, []model.Attribute{},
		[]model.Method{newMethod},
	)
	require.NoError(t, err)

	design := model.NewDesign(
		"Nonzero Sum", "Appkit is an application development kit, part of the Nonzero Sum Stack.",
		"appkit",
		[]model.Component{app, command, cli, commandFactory},
	)
	return design
}
