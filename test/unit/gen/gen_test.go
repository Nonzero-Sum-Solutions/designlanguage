package gen

import (
	"testing"

	"github.com/Nonzero-Sum-Solutions/designlanguage/gen"
	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/test/unit"
	"github.com/stretchr/testify/require"
)

func TestParseDLMFilePath1(t *testing.T) {
	ns, fn, err := gen.ParseDLMFilePath("/foo/bar/projects/baz/documentation/design", "/foo/bar/projects/baz/documentation/design/Thing1.Design.md")

	unit.RequireNoError(t, err)
	require.Equal(t, "", ns)
	require.Equal(t, "Thing1.Design.md", fn)
}

func TestParseDLMFilePath2(t *testing.T) {
	ns, fn, err := gen.ParseDLMFilePath("/foo/bar/projects/baz/documentation/design", "/foo/bar/projects/baz/documentation/design/universe/Thing1.Design.md")

	unit.RequireNoError(t, err)
	require.Equal(t, "universe", ns)
	require.Equal(t, "Thing1.Design.md", fn)
}

func TestParseDLMFilePath3(t *testing.T) {
	ns, fn, err := gen.ParseDLMFilePath("/foo/bar/projects/baz/documentation/design", "/foo/bar/projects/baz/documentation/design/universe/sol/Earth.Design.md")

	unit.RequireNoError(t, err)
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
