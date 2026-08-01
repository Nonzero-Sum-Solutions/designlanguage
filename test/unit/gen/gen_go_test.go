package gen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattmunz/designlanguage/gen"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/test/unit"
)

func TestParseDLMFilePath1(t *testing.T) {
	ns, fn, err := gen.ParseDLMFilePath("/foo/bar/projects/baz/documentation/design", "/foo/bar/projects/baz/documentation/design/Thing1.Design.md")

	require.NoError(t, err)
	require.Equal(t, "", ns)
	require.Equal(t, "Thing1.Design.md", fn)
}

func TestParseDLMFilePath2(t *testing.T) {
	ns, fn, err := gen.ParseDLMFilePath("/foo/bar/projects/baz/documentation/design", "/foo/bar/projects/baz/documentation/design/universe/Thing1.Design.md")

	require.NoError(t, err)
	require.Equal(t, "universe", ns)
	require.Equal(t, "Thing1.Design.md", fn)
}

func TestParseDLMFilePath3(t *testing.T) {
	ns, fn, err := gen.ParseDLMFilePath("/foo/bar/projects/baz/documentation/design", "/foo/bar/projects/baz/documentation/design/universe/sol/Earth.Design.md")

	require.NoError(t, err)
	require.Equal(t, "universe/sol", ns)
	require.Equal(t, "Earth.Design.md", fn)
}

func TestRenderDesignSummary(t *testing.T) {
	spinner := unit.NewSpinner(t)
	plane := unit.NewPlane(t)
	rectangle := unit.NewRectangle(t)

	design := model.NewDesign("", "", "Geometry", []model.Component{plane, rectangle, spinner})
	source := gen.RenderDesignSummary(design)
	require.Equal(t, "Namespace: Geometry, All Components: 3, Base Components: 1, Entities: 1, Objects: 1", source)
}

func TestRenderEmptyComponent(t *testing.T) {
	expectedSource := "type Plane interface{}"

	plane := unit.NewPlane(t)
	src, err := gen.RenderComponentSource(plane)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderEntity(t *testing.T) {
	expectedSource := `type Rectangle interface {
	Length() int
	Width() int
}`

	rectangle := unit.NewRectangle(t)
	src, err := gen.RenderEntitySource(rectangle)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)

}

// TestTypeAlias verifies that type names are rendered to their correlates in Go types when needed.
func TestTypeAlias(t *testing.T) {
	expectedSource := `type LabelledValue interface {
	Label() string
	Value() int
}`
	labelledValue := unit.NewLabelledValue(t)
	src, err := gen.RenderEntitySource(labelledValue)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderObject(t *testing.T) {
	expectedSource := `type Spinner interface {
	Radius() int
	Spin(velocity int, duration int)
}`

	src, err := gen.RenderObjectSource(unit.NewSpinner(t))

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderPersonRepository(t *testing.T) {
	expectedSource := `type PersonRepository interface {
	Add(person Person)
	Get(name string) (person Person)
}`

	src, err := gen.RenderObjectSource(unit.NewPersonRepo(t))

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderDesignGoSourceGeometry(t *testing.T) {
	expectedSource := `package geometry

type Plane interface{}

type Rectangle interface {
	Length() int
	Width() int
}

type Spinner interface {
	Radius() int
	Spin(velocity int, duration int)
}
`

	spinner := unit.NewSpinner(t)
	plane := unit.NewPlane(t)
	rectangle := unit.NewRectangle(t)

	design := model.NewDesign("", "", "geometry", []model.Component{plane, rectangle, spinner})

	src, err := gen.RenderDesignGoSource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestMainRenderEnumGo(t *testing.T) {
	expectedSource := `type Shape int

const (
	Ellipse Shape = iota
	Torus
	Star
	Rhombus
)`

	enumExpr, err := model.NewEnum("Shape", "Ellipse", "Torus", "Star", "Rhombus")
	require.NoError(t, err)

	src, err := gen.RenderEnumSource(enumExpr)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderDesignGoSourceSupertype(t *testing.T) {
	expectedSource := `package shapes

type Color int

const (
	Red Color = iota
	Green
	Blue
)

type Rectangle interface {
	Length() int
	Width() int
}

type Square interface {
	Rectangle

	Area() int
}

type SpinningSquare interface {
	Square

	Spin(velocity int, duration int)
}
`

	rectangle := unit.NewRectangle(t)
	square := unit.NewSquare(t)
	spinningSquare := unit.NewSpinningSquare(t)
	enum1 := unit.NewColorEnum(t)
	design := model.NewDesign("", "", "shapes", []model.Component{enum1, rectangle, square, spinningSquare})

	src, err := gen.RenderDesignGoSource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderEmptySubclass(t *testing.T) {
	expectedSource := `package pyast

type Object interface{}

type Statement interface {
	Node
}

type Position interface {
	Line() int
	ColOffset() int
}

type Node interface {
	Object

	Position() Position
}
`
	o, err := model.NewComponent("Object", "", nil)
	require.NoError(t, err)

	p, err := model.NewEntity(
		"Position", "", nil,
		[]model.Attribute{unit.NewAttr(t, "Line", "", "Int", false), unit.NewAttr(t, "ColOffset", "", "Int", false)},
	)
	require.NoError(t, err)

	n, err := model.NewEntity(
		"Node", "", unit.NewType(t, "Object"),
		[]model.Attribute{unit.NewAttr(t, "Position", "", "Position", false)},
	)
	require.NoError(t, err)

	s, err := model.NewComponent("Statement", "", unit.NewType(t, "Node"))
	require.NoError(t, err)

	design := model.NewDesign("", "", "pyast", []model.Component{o, p, n, s})

	src, err := gen.RenderDesignGoSource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderDesignGoSourceAppkit(t *testing.T) {
	expectedSource := `// Appkit is an application development kit, part of the Nonzero Sum Stack.
// Author: Nonzero Sum
package appkit

import (
	log "github.com/go-kit/kit/log"
	cobra "github.com/spf13/cobra"
)

// A software application.
type App interface {
	ID() string
	// Semver preferred.
	Version() string
	ConfigName() string
}

type Command interface {
	Execute()
}

// Command Line Interface.
type CLI interface {
	Command

	AppID() string
	Name() string
	ShortDescription() string
	NewRootCommand() (cmd *cobra.Command, err error)
	SetLogger(logger log.Logger)
}

type CommandFactory interface {
	New() (cmd Command)
}
`

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

	src, err := gen.RenderDesignGoSource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}
