package parser

import (
	"testing"

	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/parser"
	"github.com/Nonzero-Sum-Solutions/designlanguage/test/unit"
	"github.com/stretchr/testify/require"
)

func TestParseMethod1(t *testing.T) {
	method, err := model.NewMethod("Foo", "", []model.Param{}, []model.Param{})
	require.NoError(t, err)
	requireParsedMethodEqual(t, "Foo ()", method)
}

func TestParseEnum(t *testing.T) {
	enum1, err := model.NewEnum("Color", "Red", "Orangle", "Yellow", "Green")
	require.NoError(t, err)
	requireParsedEnumEqual(t, "Color = {Red, Orangle, Yellow, Green}", enum1)
}

func TestParseMethod2(t *testing.T) {
	params := []model.Param{model.NewParam("Bar", unit.NewType(t, "Baz"))}
	method, err := model.NewMethod("Foo", "", params, []model.Param{})
	require.NoError(t, err)
	requireParsedMethodEqual(t, "Foo (Bar Baz)", method)
}

func TestParseMethod3(t *testing.T) {
	method, err := model.NewMethod("Foo", "", []model.Param{}, []model.Param{})
	require.NoError(t, err)
	requireParsedMethodEqual(t, "Foo () -> ()", method)
}

func TestParseMethod4(t *testing.T) {
	params := []model.Param{model.NewParam("Bar", unit.NewType(t, "Baz"))}
	returnVals := []model.Param{model.NewParam("This", unit.NewType(t, "That"))}
	method, err := model.NewMethod("Foo", "", params, returnVals)
	require.NoError(t, err)
	requireParsedMethodEqual(t, "Foo (Bar Baz) -> (This That)", method)
}

func TestParseMethod5(t *testing.T) {
	params := []model.Param{model.NewParam("Bar", unit.NewType(t, "Baz"))}
	returnVals := []model.Param{model.NewParam("This", unit.NewType(t, "That"))}
	method, err := model.NewMethod("Foo", "A comment!", params, returnVals)
	require.NoError(t, err)
	requireParsedMethodEqual(t, "Foo (Bar Baz) -> (This That) -- A comment!", method)
}

func TestParseMethod6(t *testing.T) {
	params := []model.Param{
		model.NewParam("Bar", unit.NewType(t, "Baz")),
		model.NewParam("This", unit.NewType(t, "That")),
	}
	returnVals := []model.Param{model.NewParam("Which", unit.NewType(t, "What"))}
	method, err := model.NewMethod("Foo", "", params, returnVals)
	require.NoError(t, err)
	requireParsedMethodEqual(t, "Foo (Bar Baz, This That) -> (Which What)", method)
}

func requireParsedMethodEqual(t *testing.T, text string, expectedMethod model.Method) {
	method, err := parser.ParseMethod(text)
	require.NoError(t, err)

	require.Equal(t, expectedMethod.Name(), method.Name())
	require.Equal(t, expectedMethod.Params(), method.Params())
	require.Equal(t, expectedMethod.ReturnVals(), method.ReturnVals())
}

func requireParsedEnumEqual(t *testing.T, text string, expected model.Enum) {
	enum1, err := parser.ParseEnum(text)
	require.NoError(t, err)

	require.Equal(t, expected.Name(), enum1.Name())
	require.Equal(t, expected.Values(), enum1.Values())
}

func TestIsComponentLine1(t *testing.T) {
	require.False(t, parser.IsComponentLine1(""), "1")
	require.False(t, parser.IsComponentLine1(" "), "2")
	require.False(t, parser.IsComponentLine1(" Foo"), "3")
	require.True(t, parser.IsComponentLine1("Foo"), "4")
	require.False(t, parser.IsComponentLine1("foo"), "5")
	require.True(t, parser.IsComponentLine1("Foo"), "6")
	require.False(t, parser.IsComponentLine1("Foo::"), "7")
	require.False(t, parser.IsComponentLine1("Foo ::"), "8")
	require.True(t, parser.IsComponentLine1("Foo :: BarBaz"), "9")
	require.False(t, parser.IsComponentLine1("Foo :: barBaz"), "10")
}
