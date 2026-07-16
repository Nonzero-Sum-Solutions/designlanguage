package gotranslator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/mattmunz/appkit/maps"
	akstrings "github.com/mattmunz/appkit/strings"
	"github.com/mattmunz/designlanguage/model"
)

var goTypeByDLType = map[string]*goType{
	"String":      newGoType("", "", "string", false),
	"Int":         newGoType("", "", "int", false),
	"Integer":     newGoType("", "", "int", false),
	"Error":       newGoType("", "", "error", false),
	"CommandImpl": newGoType("*", packageNamesByAliasForImport["cobra"], "Command", false),
	"Logger":      newGoType("", packageNamesByAliasForImport["log"], "Logger", false),
}

type goType struct {
	prefix, packageName, name string
	isArray                   bool
}

// packageName is the fully qualified package name.
func newGoType(prefix, packageName, name string, isArray bool) *goType {
	return &goType{prefix, packageName, name, isArray}
}

var aliasesByPackageNameForImport = map[string]string{
	"github.com/go-kit/kit/log": "log",
	"github.com/spf13/cobra":    "cobra",
}

var packageNamesByAliasForImport = maps.InvertMap(aliasesByPackageNameForImport)

func TranslateDesign(design model.Design) (*jen.File, string, error) {
	packageName, err := getPackageName(design.Namespace())
	if err != nil {
		return nil, "", err
	}

	outFile := jen.NewFile(packageName)
	if c := design.Comment(); c != "" {
		outFile.PackageComment("// " + c)
	}

	outFile.ImportNames(packageNamesByAliasForImport)

	if a := design.Author(); a != "" {
		outFile.PackageComment("// Author: " + a)
	}

	stmt := outFile.Empty()

	for _, component := range design.Enums() {
		stmt = AddEnum(stmt, component)
	}

	for _, component := range design.BaseComponents() {
		stmt = AddBaseComponent(stmt, component)
	}

	for _, entity := range design.Entities() {
		stmt = AddEntity(stmt, entity)
	}

	for _, obj := range design.Objects() {
		stmt = AddObject(stmt, obj)
	}
	return outFile, "", nil
}

func getPackageName(namespace string) (string, error) {
	tokens := strings.Split(namespace, "/")

	if len(tokens) < 1 {
		return "", errors.New("Invalid namespace: " + namespace)
	}
	name := tokens[len(tokens)-1]
	if len(name) < 3 {
		return "", fmt.Errorf("Package name too short: [%s]", name)
	}

	return name, nil
}

func AddBaseComponent(stmt *jen.Statement, component model.Component) *jen.Statement {
	exp := stmt
	if component.Comment() != "" {
		exp = stmt.Comment("// " + component.Comment())
	}

	if component.Supertype() == nil {
		return exp.Type().Id(component.Name()).Interface().Line().Line()
	}

	return exp.Type().Id(component.Name()).InterfaceFunc(func(g *jen.Group) {
		addSupertype(g, component.Supertype())
	}).Line().Line()
}

func AddEnum(stmt *jen.Statement, enumExpr model.Enum) *jen.Statement {
	typeName := enumExpr.Name()

	exp := stmt.Type().Id(typeName).Int().Line()

	return exp.Const().DefsFunc(func(g *jen.Group) {
		for i, v := range enumExpr.Values() {
			if i < 1 {
				g.Id(v).Id(typeName).Op("=").Iota()
				continue
			}
			g.Id(v)
		}
	}).Line().Line()
}

func AddEntity(stmt *jen.Statement, entity model.Entity) *jen.Statement {
	exp := stmt
	if entity.Comment() != "" {
		exp = addComment(stmt, entity.Comment())
	}

	return exp.Type().Id(entity.Name()).InterfaceFunc(func(g *jen.Group) {
		if entity.Supertype() != nil {
			addSupertype(g, entity.Supertype())
		}

		for _, attr := range entity.Attributes() {
			addAttribute(g, attr)
		}
	}).Line().Line()
}

func AddObject(stmt *jen.Statement, obj model.Object) *jen.Statement {
	exp := stmt
	if obj.Comment() != "" {
		exp = addComment(stmt, obj.Comment())
	}

	return exp.Type().Id(obj.Name()).InterfaceFunc(func(g *jen.Group) {
		if obj.Supertype() != nil {
			addSupertype(g, obj.Supertype())
		}

		for _, attr := range obj.Attributes() {
			addAttribute(g, attr)
		}

		for _, method := range obj.Methods() {
			addMethod(g, method)
		}
	}).Line().Line()
}

func addComment(stmt *jen.Statement, text string) *jen.Statement {
	return stmt.Comment(fmt.Sprintf("// %s", text)).Line()
}

func addCommentToGroup(g *jen.Group, text string) *jen.Statement {
	return g.Comment(fmt.Sprintf("// %s", text))
}

func addSupertype(g *jen.Group, supertype model.Type) error {
	stmt, err := addJenForType(supertype, g)
	if err != nil {
		return err
	}

	stmt.Line()
	return nil
}

func addJenForType(t model.Type, g *jen.Group) (*jen.Statement, error) {
	goType := getGoType(t)
	packageName := goType.packageName
	name := goType.name
	op := goType.prefix

	var stmt *jen.Statement = g.Empty()
	if goType.isArray {
		stmt = g.Op("[]")
	}

	if op != "" {
		if op != "*" {
			return nil, fmt.Errorf("Unexpected prefix: %s", op)
		}

		if packageName != "" {
			return stmt.Op("*").Qual(packageName, name), nil
		}

		return stmt.Op("*").Id(name), nil
	}

	if packageName != "" {
		return stmt.Qual(packageName, name), nil
	}

	return stmt.Id(name), nil
}

func appendJenForType(type1 model.Type, stmt *jen.Statement) (*jen.Statement, error) {
	return addGoType(getGoType(type1), stmt)
}

// Add and attribute to the group.
// Attributes are special in that they have a return type but no return value name.
func addAttribute(g *jen.Group, attr model.Attribute) error {
	if attr.Comment() != "" {
		addCommentToGroup(g, attr.Comment())
	}

	stmt := g.Id(attr.Name()).Params()
	_, err := appendJenForType(attr.Type(), stmt)
	return err
}

func getGoType(type1 model.Type) *goType {
	if goType, exists := goTypeByDLType[type1.Name()]; exists {
		return newGoType(goType.prefix, goType.packageName, goType.name, type1.IsArray())
	}
	return newGoType("", "", type1.Name(), type1.IsArray())
}

func addMethod(g *jen.Group, method model.Method) error {
	next := g.Id(method.Name()).ParamsFunc(func(g *jen.Group) {
		for _, param := range method.Params() {
			addParamToGroup(g, param)
		}
	})

	if len(method.ReturnVals()) < 1 {
		return errors.New("No return val specified for method")
	}

	next.ParamsFunc(func(g *jen.Group) {
		for _, param := range method.ReturnVals() {
			addParamToGroup(g, param)
		}
	})

	return nil
}

func addParamToGroup(g *jen.Group, param model.Param) {
	stmt := g.Id(akstrings.LowercaseFirst(param.Name()))
	appendJenForType(param.Type(), stmt)
}

func addGoType(goType *goType, stmt *jen.Statement) (*jen.Statement, error) {
	packageName := goType.packageName
	name := goType.name
	op := goType.prefix

	newStmt := stmt

	if goType.isArray {
		newStmt = newStmt.Op("[]")
	}
	if op != "" {
		if op != "*" {
			return nil, fmt.Errorf("Unexpected prefix: %s", op)
		}

		if packageName != "" {
			return newStmt.Op("*").Qual(packageName, name), nil
		}

		return newStmt.Op("*").Id(name), nil
	}

	if packageName != "" {
		return newStmt.Qual(packageName, name), nil
	}

	return newStmt.Id(name), nil
}
