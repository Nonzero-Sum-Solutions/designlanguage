package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObject(t *testing.T) {
	name := "Shape"
	component, err := model.NewObject(name, "", nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, component)
	require.Equal(t, name, component.Name())
}

func TestLoad1(t *testing.T) {
	design := getDesign(t, "Test.1.nzsd.txt", "ns1", 1, 1, 2)

	for _, entity := range design.Entities() {
		switch entity.Name() {
		case "User":
			require.Len(t, len(entity.Attributes()), 0, "Wrong attribute count.")
			continue
		case "Console":
			require.Len(t, entity.Attributes(), 2, "Wrong entity count.")
			assertAttributeEqual(t, entity.Attributes()[0], "Version", "", "String", false)
			assertAttributeEqual(t, entity.Attributes()[1], "CommandHandlers", "", "CommandHandler", true)
			continue
		default:
			require.Fail(t, "Unknown entity: %s", entity.Name())
		}
	}

	for _, object := range design.Objects() {
		switch object.Name() {
		case "CommandHandler":
			assert.Equal(t, 1, len(object.Attributes()))
			assertAttributeEqual(t, object.Attributes()[0], "Description", "", "String", false)

			assert.Equal(t, 2, len(object.Methods()))

			for _, method := range object.Methods() {
				methodName := method.Name()
				switch methodName {
				case "Interpret":
					require.Equal(t, 1, len(method.Params()))
					assertParamEqual(t, method.Params()[0], "Command", "String", false)
					require.Equal(t, 1, len(method.ReturnVals()))
					assertParamEqual(t, method.ReturnVals()[0], "Response", "String", false)
					continue
				case "GetConfigs":
					require.Equal(t, 1, len(method.Params()))
					assertParamEqual(t, method.Params()[0], "Time", "DateTime", false)
					require.Equal(t, 1, len(method.ReturnVals()))
					assertParamEqual(t, method.ReturnVals()[0], "Configs", "String", true)
					continue
				default:
					require.Fail(t, "Reached unknown method.", "Unknown method: [%s]", methodName)
				}
			}

		case "Browser":
			require.Equal(t, 0, len(object.Attributes()))
			assert.Equal(t, 1, len(object.Methods()))
			method := object.Methods()[0]

			require.Equal(t, "GetElement", method.Name())
			require.Equal(t, 1, len(method.Params()))
			assertParamEqual(t, method.Params()[0], "ID", "String", false)
			require.Equal(t, 1, len(method.ReturnVals()))
			assertParamEqual(t, method.ReturnVals()[0], "Element", "Element", false)
		default:
			require.Fail(t, fmt.Sprintf("Unknown object: %s", object.Name()))
		}
	}
}

func TestLoad3(t *testing.T) {
	design := getDesign(t, "Test.3.nzsd.txt", "ns3", 0, 0, 1)

	obj := design.Objects()[0]

	require.Equal(t, "Runner", obj.Name())
	assert.Equal(t, 1, len(obj.Methods()))

	method := obj.Methods()[0]

	require.Equal(t, "Run", method.Name())
	require.Equal(t, 1, len(method.Params()))
	require.Equal(t, 0, len(method.ReturnVals()))

	assertParamEqual(t, method.Params()[0], "Args", "String", false)
}

func TestLoadAppKit(t *testing.T) {
	design := getDesign(t, "project1/documentation/design/appkit/Appkit.nzsd.txt", "appkit", 0, 1, 3)

	require.Equal(t, "Nonzero Sum", design.Author())
	require.Equal(t, "Appkit is an application development kit, part of the Nonzero Sum Stack.", design.Comment())

	require.Len(t, design.Entities(), 1)

	entity := design.Entities()[0]
	require.Equal(t, "A software application.", entity.Comment())
	require.Len(t, entity.Attributes(), 3)

	attr2 := entity.Attributes()[1]
	requireNameEquals(t, "Version", attr2)
	require.Equal(t, "Semver preferred.", attr2.Comment())

	obj1 := design.Objects()[0]
	requireNameEquals(t, "Command", obj1)
	require.Len(t, obj1.Methods(), 1)
	method1 := obj1.Methods()[0]
	requireNameEquals(t, "Execute", method1)
	require.Len(t, method1.Params(), 0)
	require.Len(t, method1.ReturnVals(), 0)

	obj2 := design.Objects()[1]
	requireNameEquals(t, "CLI", obj2)
	require.Len(t, obj2.Methods(), 2)

	method21 := obj2.Methods()[0]
	requireNameEquals(t, "NewRootCommand", method21)
	require.Equal(t, 0, len(method21.Params()))
	require.Equal(t, 2, len(method21.ReturnVals()))
	requireParamEquals(t, method21.ReturnVals()[0], "Cmd", "CommandImpl", false)

	method22 := obj2.Methods()[1]
	requireNameEquals(t, "SetLogger", method22)
	require.Equal(t, 1, len(method22.Params()))
	requireParamEquals(t, method22.Params()[0], "Logger", "Logger", false)
}

func getDesign(t *testing.T, testFileSubpath string, namespace string, expectedBaseComponentsLength int, expectedEntitiesLength int, expectedObjectsLength int) model.Design {
	absPath, err := filepath.Abs(
		filepath.Join(getTestDataPath(t), testFileSubpath),
	)
	require.NoError(t, err)

	designParser := parser.NewParser()

	require.FileExists(t, absPath)

	design, err := designParser.Parse(absPath, namespace)
	require.NoError(t, err)

	require.Len(t, design.BaseComponents(), expectedBaseComponentsLength, "Wrong base component count.")
	require.Len(t, design.Entities(), expectedEntitiesLength, "Wrong entity count.")
	require.Len(t, design.Objects(), expectedObjectsLength, "Wrong object count.")
	return design
}

func assertAttributeEqual(t *testing.T, attr model.Attribute, expectedName, expectedComment, expectedType string, expectedIsArray bool) {
	require.Equal(t, expectedName, attr.Name())
	require.Equal(t, expectedComment, attr.Comment())
	require.Equal(t, expectedType, attr.Type().Name())
	require.Equal(t, expectedIsArray, attr.Type().IsArray())
}

func requireNameEquals(t *testing.T, expectedName string, named model.Named) {
	require.Equal(t, expectedName, named.Name())
}

func assertParamEqual(t *testing.T, param model.Param, expectedName string, expectedType string, expectedIsArray bool) {
	require.Equal(t, expectedName, param.Name())
	require.Equal(t, expectedType, param.Type().Name())
	require.Equal(t, expectedIsArray, param.Type().IsArray())
}

func getTestDataPath(t *testing.T) string {
	value := os.Getenv("TEST_DATA_PATH")
	if value == "" {
		require.Fail(t, "TEST_DATA_PATH is not set")
	}
	return value
}

func requireParamEquals(t *testing.T, param model.Param, expectedName, expectedTypeName string, expectedIsArray bool) {
	require.Equal(t, expectedName, param.Name())
	require.Equal(t, expectedTypeName, param.Type().Name())
	require.Equal(t, expectedIsArray, param.Type().IsArray())
}
