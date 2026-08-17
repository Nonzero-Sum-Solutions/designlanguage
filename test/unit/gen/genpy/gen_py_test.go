package genpy

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Nonzero-Sum-Solutions/designlanguage/gen/genpy"
	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/test/unit"
	"github.com/Nonzero-Sum-Solutions/designlanguage/translator/pytranslator"
	"github.com/Nonzero-Sum-Solutions/designlanguage/writer/pywriter"

	pyast "github.com/Nonzero-Sum-Solutions/designlanguage/model/gen/pyast"
	pyastm "github.com/Nonzero-Sum-Solutions/designlanguage/model/pyast"
)

func TestRenderPyEmptyComponent(t *testing.T) {
	expectedSource := `class Plane(Protocol):
	...

`

	plane := unit.NewPlane(t)
	src, err := genpy.RenderComponentPySource(plane)

	unit.RequireNoError(t, err)
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
	src, err := genpy.RenderEntityPySource(rectangle)

	unit.RequireNoError(t, err)
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
	src, err := genpy.RenderEntityPySource(labelledValue)

	unit.RequireNoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderPyObject(t *testing.T) {
	expectedSource := `class Spinner(Protocol):
	def radius() -> int:
		...

	def spin(velocity: int, duration: int):
		...

`

	src, err := genpy.RenderObjectPySource(unit.NewSpinner(t))

	unit.RequireNoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestNewMethodDef(t *testing.T) {
	method1 := unit.NewSpinner(t).Methods()[0]
	methodDef, err := pytranslator.NewMethodDef(method1)
	unit.RequireNoError(t, err)

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

	text, err := pywriter.RenderArguments(methodDef.Args(), "")
	unit.RequireNoError(t, err)

	require.Equal(t, "(velocity: int, duration: int)", text)
}

func TestPyRenderPersonRepository(t *testing.T) {
	expectedSource := `class PersonRepository(Protocol):
	def add(person: Person):
		...

	def get(name: str) -> Person:
		...

`

	src, err := genpy.RenderObjectPySource(unit.NewPersonRepo(t))

	unit.RequireNoError(t, err)
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

	src, err := genpy.RenderDesignPySource(design)
	unit.RequireNoError(t, err)
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

	src, err := genpy.RenderDesignPySource(design)
	unit.RequireNoError(t, err)
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

	def new_root_command() -> CommandImpl:
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

	design := unit.NewAppkitDesign(t)

	src, err := genpy.RenderDesignPySource(design)
	unit.RequireNoError(t, err)
	require.Equal(t, expectedSource, src)
}
