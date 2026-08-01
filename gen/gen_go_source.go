package gen

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go/format"

	"github.com/dave/jennifer/jen"
	klog "github.com/go-kit/kit/log"
	akstrings "github.com/mattmunz/appkit/strings"

	"github.com/mattmunz/appkit/misc"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/parser"
)

const formatSourceCode = true

var goTypeByDLType = map[string]*goType{
	"String":  newGoType("", "", "string", false),
	"Int":     newGoType("", "", "int", false),
	"Integer": newGoType("", "", "int", false),
	"Error":   newGoType("", "", "error", false),
	// TODO Extract to config file
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

// TODO Extract to config file
var aliasesByPackageNameForImport = map[string]string{
	"github.com/go-kit/kit/log": "log",
	"github.com/spf13/cobra":    "cobra",
}

var packageNamesByAliasForImport = InvertMap(aliasesByPackageNameForImport)

// TODO Move to appkit
// InvertMap swaps keys and values.
// The value type (V) must be comparable to be used as a map key.
func InvertMap[K comparable, V comparable](original map[K]V) map[V]K {
	inverted := make(map[V]K, len(original))
	for k, v := range original {
		inverted[v] = k
	}
	return inverted
}

func RenderDesignSummary(design model.Design) string {
	return fmt.Sprintf(
		"Namespace: %s, All Components: %d, Base Components: %d, Entities: %d, Objects: %d",
		design.Namespace(), len(design.AllComponents()), len(design.BaseComponents()),
		len(design.Entities()), len(design.Objects()),
	)
}

func RenderDesignGoSource(design model.Design) (string, error) {
	packageName, err := getPackageName(design.Namespace())
	if err != nil {
		return "", err
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
		stmt = addEnum(stmt, component)
	}

	for _, component := range design.BaseComponents() {
		stmt = addBaseComponent(stmt, component)
	}

	for _, entity := range design.Entities() {
		stmt = addEntity(stmt, entity)
	}

	for _, obj := range design.Objects() {
		stmt = addObject(stmt, obj)
	}

	return renderGoFile(outFile)
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

func RenderEnumSource(enumExpr model.Enum) (string, error) {
	stmt := jen.Empty()
	addEnum(stmt, enumExpr)
	return renderGo(stmt)
}

func RenderComponentSource(component model.Component) (string, error) {
	stmt := jen.Empty()
	addBaseComponent(stmt, component)
	return renderGo(stmt)
}

// TODO SOme shared code here that could be refactored together
func addBaseComponent(stmt *jen.Statement, component model.Component) *jen.Statement {
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

func addEnum(stmt *jen.Statement, enumExpr model.Enum) *jen.Statement {
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

func RenderEntitySource(entity model.Entity) (string, error) {
	stmt := jen.Empty()
	addEntity(stmt, entity)
	return renderGo(stmt)
}

func addEntity(stmt *jen.Statement, entity model.Entity) *jen.Statement {
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

func addObject(stmt *jen.Statement, obj model.Object) *jen.Statement {
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
	return renderGoType(getGoType(type1), stmt)
}

func renderGoType(goType *goType, stmt *jen.Statement) (*jen.Statement, error) {
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

// Only for testing.
func RenderObjectSource(obj model.Object) (string, error) {
	stmt := jen.Empty()
	addObject(stmt, obj)
	return renderGo(stmt)
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

func renderGoFile(f *jen.File) (string, error) {
	// TODO f.GoString is only for testing!
	text := f.GoString()
	if !formatSourceCode {
		return text, nil
	}
	src, err := format.Source([]byte(text))
	if err != nil {
		return "", fmt.Errorf("Couldn't format source code 1: %w", err)
	}
	return string(src), err
}

func renderGo(statement *jen.Statement) (string, error) {
	if !formatSourceCode {
		return statement.GoString(), nil
	}

	src, err := format.Source([]byte(strings.TrimSpace(statement.GoString())))
	if err != nil {
		return "", fmt.Errorf("Couldn't format source code 2: %w", err)
	}
	return string(src), nil
}

func GenerateGoSourceForDL(projectDir string, logger klog.Logger, dryRun bool) error {
	projectDirPath, err := filepath.Abs(projectDir)
	if err != nil {
		return err
	}
	designPath := filepath.Join(projectDirPath, "documentation", "design")

	designDirInfo, err := os.Stat(designPath)
	if err != nil {
		return err
	}

	if !designDirInfo.IsDir() {
		return fmt.Errorf("File is not dir: %s", designPath)
	}

	parsedDesigns := NewDesignList()

	misc.LogMessage(logger, fmt.Sprintf("Walking the path %q", designPath))

	designParser := parser.NewParser()

	err = filepath.Walk(designPath, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".nzsd.txt") {
			return nil
		}

		return HandleDLMFile(logger, designParser, parsedDesigns, designPath, path, info, dryRun, projectDirPath, err)
	})

	if err != nil {

		return fmt.Errorf("Error walking the path %q: %w", designPath, err)
	}

	misc.LogMessage(logger, "Generating code...")

	for _, design := range parsedDesigns.All() {
		misc.LogMessage(logger, fmt.Sprintf("Design summary:\n%s", RenderDesignSummary(design)))

		designSource, err := RenderDesignGoSource(design)

		if err != nil {
			misc.LogMessage(logger, fmt.Sprintf("Error rendering design source: %v", err))
			continue
		}

		misc.LogMessage(logger, fmt.Sprintf("Design source: %s", designSource))

		if dryRun {
			continue
		}

		dirPath := filepath.Join(projectDirPath, "model", "gen", design.Namespace())
		filePath := filepath.Join(dirPath, fmt.Sprintf("%s.go", design.Namespace()))

		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if err = os.WriteFile(filePath, []byte(designSource), 0644); err != nil {
			return err
		}

		misc.LogMessage(logger, fmt.Sprintf("Wrote file: %s", filePath))
	}
	return nil
}
