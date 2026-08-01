package gen

// todo this is messed up. MAke each statement or expression end with the correct amount of trailing space. This will depend on container context
import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattmunz/designlanguage/gen"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/test/unit"

	pyast "github.com/mattmunz/designlanguage/model/gen/pyast"
	pyastm "github.com/mattmunz/designlanguage/model/pyast"
)

func TestRenderPyEmptyComponent(t *testing.T) {
	expectedSource := `class Plane(Protocol):
	...

`

	plane := unit.NewPlane(t)
	src, err := gen.RenderComponentPySource(plane)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderPyEntity(t *testing.T) {
	expectedSource := `class Rectangle(Protocol):
	def length() -> int:
		...

	def width() -> int:
		...

`

	rectangle := unit.NewRectangle(t)
	src, err := gen.RenderEntityPySource(rectangle)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)

}

// TestGoTypeAlias verifies that type names are rendered to their correlates in Go types when needed.
func TestTypePyAlias(t *testing.T) {
	expectedSource := `class LabelledValue(Protocol):
	def label() -> str:
		...

	def value() -> int:
		...

`

	labelledValue := unit.NewLabelledValue(t)
	src, err := gen.RenderEntityPySource(labelledValue)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderPyObject(t *testing.T) {
	expectedSource := `class Spinner(Protocol):
	def radius() -> int:
		...

	def spin(velocity: int, duration: int):
		...

`

	src, err := gen.RenderObjectPySource(unit.NewSpinner(t))

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestNewMethodDef(t *testing.T) {
	method1 := unit.NewSpinner(t).Methods()[0]
	methodDef, err := gen.NewMethodDef(method1)
	require.NoError(t, err)

	require.Equal(t, "spin", methodDef.Name().Value())
	args := methodDef.Args().Args()
	require.Len(t, args, len(method1.Params()))
	require.Equal(t, "velocity", args[0].Arg().Value())
	annotation1 := args[0].Annotation()

	name := pyastm.NewName("foo", pyast.Load)
	if reflect.TypeOf(annotation1) != reflect.TypeOf(name) {
		t.Errorf("got type %v, want %v", reflect.TypeOf(annotation1), reflect.TypeOf(name))
	}

	n, ok := annotation1.(pyast.Name)
	require.True(t, ok)

	require.Equal(t, n.Id().Value(), "int")

	text, err := gen.RenderArguments(methodDef.Args(), "")
	require.NoError(t, err)

	require.Equal(t, "(velocity: int, duration: int)", text)
}

func TestPyRenderPersonRepository(t *testing.T) {
	expectedSource := `class PersonRepository(Protocol):
	def add(person: Person):
		...

	def get(name: str) -> Person:
		...

`

	src, err := gen.RenderObjectPySource(unit.NewPersonRepo(t))

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestPyRenderDesignGoSourceGeometry(t *testing.T) {
	expectedSource := `from typing import Protocol

class Plane(Protocol):
	...

class Rectangle(Protocol):
	def length() -> int:
		...

	def width() -> int:
		...

class Spinner(Protocol):
	def radius() -> int:
		...

	def spin(velocity: int, duration: int):
		...

`

	spinner := unit.NewSpinner(t)
	plane := unit.NewPlane(t)
	rectangle := unit.NewRectangle(t)

	design := model.NewDesign("", "", "geometry", []model.Component{plane, rectangle, spinner})

	src, err := gen.RenderDesignPySource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderDesignPySourceSupertype(t *testing.T) {
	expectedSource := `from typing import Protocol

class Plane(Protocol):
	...

class Rectangle(Protocol):
	def length() -> int:
		...

	def width() -> int:
		...

class Square(Rectangle):
	def area() -> int:
		...

class SpinningSquare(Square):
	def spin(velocity: int, duration: int):
		...

`
	plane := unit.NewPlane(t)
	rectangle := unit.NewRectangle(t)
	square := unit.NewSquare(t)
	spinningSquare := unit.NewSpinningSquare(t)
	design := model.NewDesign("", "", "shapes", []model.Component{plane, rectangle, square, spinningSquare})

	src, err := gen.RenderDesignPySource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestGoRenderDesignPySourceAppkit(t *testing.T) {
	// TODO This needds tweaking
	expectedSource := `# Appkit is an application development kit, part of the Nonzero Sum Stack.
from typing import Protocol

class App(Protocol):
	"""A software application."""
	def id() -> str:
		...

	def version() -> str:
		"""Semver preferred."""
		...

	def config_name() -> str:
		...

class Command(Protocol):
	def execute():
		...

class CLI(Command):
	"""Command Line Interface."""
	def app_id() -> str:
		...

	def name() -> str:
		...

	def short_description() -> str:
		...

	def new_root_command() -> Command:
		"""
		Raises:
			Exception
		"""
		...

	def set_logger(logger: Logger):
		...

class CommandFactory(Protocol):
	def new() -> Command:
		...

`

	// TODO Refactor with go test via unit
	appAttrs := []model.Attribute{
		unit.NewAttr(t, "ID", "", "String", false),
		unit.NewAttr(t, "Version", "Semver preferred.", "String", false),
		unit.NewAttr(t, "ConfigName", "", "String", false),
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
		unit.NewAttr(t, "AppID", "", "String", false),
		unit.NewAttr(t, "Name", "", "String", false),
		unit.NewAttr(t, "ShortDescription", "", "String", false),
	}

	nrReturnVals := []model.Param{
		model.NewParam("Cmd", unit.NewType(t, "CommandImpl")),
		model.NewParam("Err", unit.NewType(t, "Error")),
	}
	nrMethod, err := model.NewMethod("NewRootCommand", "", []model.Param{}, nrReturnVals)
	require.NoError(t, err)

	slMethod, err := model.NewMethod("SetLogger", "", []model.Param{model.NewParam("Logger", unit.NewType(t, "Logger"))}, []model.Param{})
	require.NoError(t, err)

	cliMethods := []model.Method{nrMethod, slMethod}
	cli, err := model.NewObject("CLI", "Command Line Interface.", unit.NewType(t, "Command"), cliAttributes, cliMethods)
	require.NoError(t, err)

	newMethod, err := model.NewMethod("New", "", []model.Param{}, []model.Param{model.NewParam("Cmd", unit.NewType(t, "Command"))})
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

	src, err := gen.RenderDesignPySource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}
