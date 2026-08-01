package gengo

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Nonzero-Sum-Solutions/designlanguage/gen/gengo"
	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/test/unit"
)

func TestRenderEmptyComponent(t *testing.T) {
	expectedSource := "type Plane interface{}"

	plane := unit.NewPlane(t)
	src, err := gengo.RenderComponentSource(plane)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderEntity(t *testing.T) {
	expectedSource := `type Rectangle interface {
	Length() int
	Width() int
}`

	rectangle := unit.NewRectangle(t)
	src, err := gengo.RenderEntitySource(rectangle)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)

}

func TestObjectAlias(t *testing.T) {
	expectedSource := `type Shape interface {
	ID() any
}`

	shape, err := model.NewEntity(
		"Shape", "", nil,
		[]model.Attribute{unit.NewAttr(t, "ID", "", "Object", false)},
	)

	src, err := gengo.RenderEntitySource(shape)

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
	src, err := gengo.RenderEntitySource(labelledValue)

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderObject(t *testing.T) {
	expectedSource := `type Spinner interface {
	Radius() int
	Spin(velocity int, duration int)
}`

	src, err := gengo.RenderObjectSource(unit.NewSpinner(t))

	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}

func TestRenderPersonRepository(t *testing.T) {
	expectedSource := `type PersonRepository interface {
	Add(person Person)
	Get(name string) (person Person)
}`

	src, err := gengo.RenderObjectSource(unit.NewPersonRepo(t))

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

	src, err := gengo.RenderDesignGoSource(design)
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

	src, err := gengo.RenderEnumSource(enumExpr)
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

	src, err := gengo.RenderDesignGoSource(design)
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

	src, err := gengo.RenderDesignGoSource(design)
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

	design := unit.NewAppkitDesign(t)

	src, err := gengo.RenderDesignGoSource(design)
	require.NoError(t, err)
	require.Equal(t, expectedSource, src)
}
