package gen

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	klog "github.com/go-kit/kit/log"

	"github.com/iancoleman/strcase"

	"github.com/mattmunz/appkit/misc"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/model/gen/pyast"
	pyastm "github.com/mattmunz/designlanguage/model/pyast"
	"github.com/mattmunz/designlanguage/parser"
)

// TODO Split files into gen model, extended model (modelx?), AST builder (adapt DL to AST), Source Writer (adapt from AST to source)

var pyTypeByDLType = map[string]*pyType{
	"String": newPyType("", "", "str", false),
	"Int":    newPyType("", "", "int", false),
	// TODO This needs to be exception, right?
	"Error": newPyType("", "", "error", false),
	// TODO Extract to config file
	"CommandImpl": newPyType("*", packageNamesByAliasForImport["cobra"], "Command", false),
	"Logger":      newPyType("", packageNamesByAliasForImport["log"], "Logger", false),
}

// TODO Rework
type pyType struct {
	op, packageName, name string
	isArray               bool
}

// TODO Rework
// packageName is the fully qualified package name.
func newPyType(prefix, packageName, name string, isArray bool) *pyType {
	return &pyType{prefix, packageName, name, isArray}
}

// TODO Extract to config file, fix for py
var pyAliasesByPackageNameForImport = map[string]string{
	"github.com/go-kit/kit/log": "log",
	"github.com/spf13/cobra":    "cobra",
}

var pyPackageNamesByAliasForImport = InvertMap(pyAliasesByPackageNameForImport)

func RenderDesignPySource(design model.Design) (string, error) {
	ast1 := pyastm.NewAST()

	if design.Comment() != "" {
		ast1.Add(pyastm.NewCommentLine(design.Comment()))
	}

	// TODO IN future will need more sophistocated import mechanism
	ast1.Add(pyastm.NewImportFrom(1, pyastm.Id("typing"), []pyast.Alias{pyastm.NewAlias("Protocol")}))

	err := addDesignPy(ast1, design)
	if err != nil {
		return "", err
	}

	text, err := renderPySource(ast1)
	if err != nil {
		return "", err
	}
	return text, nil
}

func addDesignPy(ast1 pyast.AST, design model.Design) error {
	for _, e := range design.Enums() {
		addEnumPy(ast1, e)
	}

	for _, c := range design.BaseComponents() {
		addBaseComponentPy(ast1, c)
	}

	for _, en := range design.Entities() {
		addEntityPy(ast1, en)
	}

	for _, o := range design.Objects() {
		addObjectPy(ast1, o)
	}

	return nil
}

// TODO Snake case this?
func getPyType(type1 model.Type) *pyType {
	if pyType, exists := pyTypeByDLType[type1.Name()]; exists {
		return newPyType(pyType.op, pyType.packageName, pyType.name, type1.IsArray())
	}
	return newPyType("", "", type1.Name(), type1.IsArray())
}

// TODO Redundant
func renderPySource(ast1 pyast.AST) (string, error) {
	return renderPy(ast1, "", "")
}

func RenderComponentPySource(component model.Component) (string, error) {
	var ast1 pyast.AST = pyastm.NewAST()
	addBaseComponentPy(ast1, component)
	return renderPy(ast1, "", "")
}

// TODO Impl TDD
func addEnumPy(ast1 pyast.AST, e model.Enum) {
	panic("TODO NYI aep")
}

func addBaseComponentPy(ast1 pyast.AST, component model.Component) {
	body := []pyast.Statement{}
	if component.Comment() != "" {
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(pyastm.NewCommentLine(component.Comment()))))
	}
	body = append(body, pyastm.NewExprStatement(pyastm.NewEllipsis()))

	ast1.Add(pyastm.NewClassDef(pyastm.Id(component.Name()), getClassDefBases(component), nil, nil, nil, body, nil))
}

func RenderEntityPySource(entity model.Entity) (string, error) {
	var ast1 pyast.AST = pyastm.NewAST()
	addEntityPy(ast1, entity)
	return renderPy(ast1, "", "")
}

/*
Docstrings: Triple-quoted strings (""" doc """) at the beginning of modules or functions are not comments.
They are parsed as regular statement nodes (ast.Expr containing an ast.Constant) and can always be retrieved
using the built-in ast.get_docstring() function.
*/
// TODO MAke a test with a doc string comment and a end of line comment. Comment should be a statement type, and maybe docstring and expression
// REfactor with addBaseComponentPy
func addEntityPy(ast1 pyast.AST, entity model.Entity) {
	// TODO Refacrot with componentPy
	body := []pyast.Statement{}
	if entity.Comment() != "" {
		comment := pyastm.NewCommentLine(entity.Comment())
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(comment)))
	}

	for _, a := range entity.Attributes() {
		body = append(body, newAttrDef(a))
	}

	ast1.Add(pyastm.NewClassDef(pyastm.Id(entity.Name()), getClassDefBases(entity), nil, nil, nil, body, nil))
}

func getClassDefBases(entity model.Component) []pyast.Expression {
	if entity.Supertype() == nil {
		baseExp := pyastm.NewName("Protocol", pyast.Load)
		return []pyast.Expression{baseExp}
	}
	return []pyast.Expression{pyastm.NewName(entity.Supertype().Name(), pyast.Load)}
}

func newAttrDef(a model.Attribute) pyast.FunctionDef {
	// TODO Refactor with similar
	body := []pyast.Statement{}
	if a.Comment() != "" {
		comment := pyastm.NewCommentLine(a.Comment())
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(comment)))
	}
	body = append(body, pyastm.NewExprStatement(pyastm.NewEllipsis()))

	// TODO refactor TODO wrong -- using return but should be naked expression
	pyType := getPyType(a.Type())
	args := pyastm.NewArguments([]pyast.Arg{})
	return pyastm.NewFunctionDef(toPyFnCase(a.Name()), args, body, toPYASTTypeExpression(pyType))
}

func toPYASTTypeExpression(pyType *pyType) pyast.Expression {
	if !pyType.isArray {
		return pyastm.NewName(pyType.name, pyast.Load)
	}
	// TODO TYhis right?
	return pyastm.NewName("List<"+pyType.name+">", pyast.Load)
}

// TODO Needs return statement
func NewMethodDef(m model.Method) (pyast.FunctionDef, error) {

	var returnsStmt pyast.Expression
	var err error

	// TODO refactor TODO wrong -- using return but should be naked expression
	returnVals := m.ReturnVals()

	raisesException := false
	// TODO Need to handle when error is a ret val. Need a doc string there
	elts := []pyast.Expression{}
	for _, v := range returnVals {
		if v.Type().Name() == "Error" {
			raisesException = true
			continue
		}

		var elt pyast.Expression
		pyType := getPyType(v.Type())
		// TODO Refactor with similar
		if !pyType.isArray {
			elt = pyastm.NewName(pyType.name, pyast.Load)
		} else {
			// TODO TYhis right?
			elt = pyastm.NewName("List<"+pyType.name+">", pyast.Load)
		}
		elts = append(elts, elt)
	}

	// TOOOD Refactor with similar
	body := []pyast.Statement{}
	if m.Comment() != "" || raisesException {
		// MAke a ddoc string
		lines := []pyast.CommentLine{}
		if m.Comment() != "" {
			comment := pyastm.NewCommentLine(m.Comment())
			lines = append(lines, comment)
		}
		if raisesException {
			lines = append(lines, pyastm.NewCommentLine("Raises:"), pyastm.NewCommentLine("\tException"))
		}
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(lines...)))
	}
	body = append(body, pyastm.NewExprStatement(pyastm.NewEllipsis()))

	if len(elts) < 1 {
		returnsStmt = nil
	} else {
		returnsStmt, err = pyastm.NewTuple(elts, pyast.Load)
		if err != nil {
			return nil, err
		}
	}

	args := []pyast.Arg{}
	for _, p := range m.Params() {
		pyType := getPyType(p.Type())
		argName := toPyArgCase(p.Name())
		arg := pyastm.NewArg(pyastm.Id(argName), toPYASTTypeExpression(pyType))
		args = append(args, arg)
	}
	args2 := pyastm.NewArguments(args)
	return pyastm.NewFunctionDef(toPyFnCase(m.Name()), args2, body, returnsStmt), nil
}

// TODO Needs to do the snake case thing
func toPyFnCase(s string) string {
	return strcase.ToSnake(s)
}

// TODO Refactor
// TODO Needs to do the snake case thing
func toPyArgCase(s string) string {
	return strcase.ToSnake(s)
}

func addObjectPy(ast1 pyast.AST, obj model.Object) error {
	// TODO Refactor with similar
	body := []pyast.Statement{}
	if obj.Comment() != "" {
		comment := pyastm.NewCommentLine(obj.Comment())
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(comment)))
	}
	if len(obj.Attributes())+len(obj.Methods()) == 0 {
		// SHould be impossible
		body = append(body, pyastm.NewExprStatement(pyastm.NewEllipsis()))
	}

	for _, a := range obj.Attributes() {
		body = append(body, newAttrDef(a))
	}

	for _, m := range obj.Methods() {
		methodDef, err := NewMethodDef(m)
		if err != nil {
			return err
		}
		body = append(body, methodDef)
	}

	ast1.Add(pyastm.NewClassDef(pyastm.Id(obj.Name()), getClassDefBases(obj), nil, nil, nil, body, nil))
	return nil
}

// Only for testing.
func RenderObjectPySource(obj model.Object) (string, error) {
	var ast1 pyast.AST = pyastm.NewAST()
	if err := addObjectPy(ast1, obj); err != nil {
		return "", err
	}
	return renderPy(ast1, "", "")
}

func renderPy(ast1 pyast.AST, linePrefix, container string) (string, error) {
	blocks := []string{}
	for _, n := range ast1.Nodes() {
		switch v := n.(type) {
		// TODO Handle enums
		case pyast.ImportFrom:
			// IF last import statement...
			blocks = append(blocks, renderImportFrom(v)+"\n\n")
		case pyast.ClassDef:
			text, err := renderClassDef(v, linePrefix)
			if err != nil {
				return "", err
			}
			blocks = append(blocks, text)
		case pyast.CommentLine:
			text := renderCommentLine(v, linePrefix, container)
			blocks = append(blocks, text)
		default:
			return "", fmt.Errorf("Unknown node type: %#v", n)
		}
	}

	return strings.Join(blocks, ""), nil
}

// TODO Handle "asname" whern needed
func renderImportFrom(i pyast.ImportFrom) string {
	nameTexts := []string{}
	for _, n := range i.Names() {
		nameTexts = append(nameTexts, n.Name().Value())
	}

	return "from " + i.Module().Value() + " import " + strings.Join(nameTexts, ", ")
}

// TODO Use stringtermplate?
func renderClassDef(c pyast.ClassDef, linePrefix string) (string, error) {
	basesText := ""
	if len(c.Bases()) > 0 {
		basesVals, err := getBaseNames(c)
		if err != nil {
			return "", err
		}
		basesText = "(" + strings.Join(basesVals, ", ") + ")"
	}
	src := linePrefix + "class " + c.Name().Value() + basesText + ":\n"

	bodyStatements := []string{}
	for _, s := range c.Body() {
		text, err := renderStatement(s, linePrefix+"\t", "classDef")
		if err != nil {
			return "", err
		}
		bodyStatements = append(bodyStatements, text)
	}
	return src + strings.Join(bodyStatements, "\n") + "\n", nil
}

func getBaseNames(c pyast.ClassDef) ([]string, error) {
	names := []string{}
	for _, b := range c.Bases() {
		switch v := b.(type) {
		case pyast.Identifier:
			names = append(names, v.Value())
		case pyast.Name:
			names = append(names, v.Id().Value())
		default:
			return nil, fmt.Errorf("Unexpected expression type: v+%", v)
		}
	}
	return names, nil
}

func renderStatement(s pyast.Statement, linePrefix, container string) (string, error) {
	switch v := s.(type) {
	case pyast.FunctionDef:
		return RenderFunctionDef(v, linePrefix, container)
	case pyast.Name:
		return renderName(v), nil
	case pyast.ExprStatement:
		return renderExpression(v.ExprValue(), linePrefix, " ", container)
	case pyast.Return:
		return renderExpression(v.ReturnValue(), linePrefix, " ", container)
	case pyast.CommentLine:
		return renderCommentLine(v, linePrefix, container), nil
	// case pyast.Tuple:
	// 	return renderTuple(v, linePrefix)
	default:
		return "", fmt.Errorf("Unknown statement type: %#v", v)
	}
}

func renderCommentLine(v pyast.CommentLine, linePrefix, container string) string {
	if container == "" {
		return linePrefix + "# " + v.Text() + "\n"
	}
	if container == "docstring1" {
		return linePrefix + v.Text()
	}
	return linePrefix + v.Text() + "\n"
}

// TODO See if can use this in more places
func renderStatements(stmts []pyast.Statement, linePrefix, separator, container string) (string, error) {
	texts := []string{}
	for _, s := range stmts {
		text, err := renderStatement(s, linePrefix, container)
		if err != nil {
			return "", err
		}
		texts = append(texts, text)
	}
	return strings.Join(texts, separator), nil
}

func renderExpression(e pyast.Expression, linePrefix, separator, container string) (string, error) {
	switch v := e.(type) {
	case pyast.Name:
		return renderName(v), nil
	case pyast.Tuple:
		return renderTuple(v, linePrefix, container)
	case pyast.Ellipsis:
		if container == "classDef" {
			return linePrefix + v.Symbol() + "\n", nil
		}
		return linePrefix + v.Symbol(), nil
	case pyast.DocString:
		return renderDocString(v, linePrefix)
	}
	return renderStatements(e.Body(), linePrefix, separator, container)
}

func renderDocString(v pyast.DocString, linePrefix string) (string, error) {
	body := v.Body()

	if len(body) < 1 {
		return "", errors.New("Empty docstring")
	}

	if len(body) == 1 {
		text, err := renderStatement(body[0], "", "docstring1")
		if err != nil {
			return "", err
		}
		return linePrefix + "\"\"\"" + text + "\"\"\"", nil
	}

	lines := []string{}
	lines = append(lines, linePrefix+"\"\"\"")
	for _, s := range body {
		text2, err := renderStatement(s, linePrefix, "docstring")
		if err != nil {
			return "", err
		}
		lines = append(lines, text2)
	}
	lines = append(lines, linePrefix+"\"\"\"")

	return strings.Join(lines, ""), nil
}

func renderName(n pyast.Name) string {
	return n.Id().Value()
}

func renderTuple(t pyast.Tuple, linePrefix, container string) (string, error) {
	expTexts := []string{}
	for _, e := range t.Elts() {
		text, err := renderExpression(e, "", " ", container)
		if err != nil {
			return "", err
		}
		expTexts = append(expTexts, linePrefix+text)
	}

	return strings.Join(expTexts, ", "), nil
}

func RenderFunctionDef(f pyast.FunctionDef, linePrefix, container string) (string, error) {
	argsText, err := RenderArguments(f.Args(), container)
	if err != nil {
		return "", nil
	}
	line1 := "def " + f.Name().Value() + argsText

	returns := f.Returns()
	if returns != nil {
		stmt1, err := renderExpression(returns, "", "", container)
		if err != nil {
			return "", err
		}
		if stmt1 == "" {
			return "", fmt.Errorf("Empty returns statement. FunctionDef: %#v, returns: %#v", f, returns)
		}
		line1 = line1 + " -> " + stmt1
	}

	text := linePrefix + line1 + ":" + "\n"

	// TODO Refactor using renderStatements method
	bodyLines := []string{}
	for _, s := range f.Body() {
		stmt2, err := renderStatement(s, linePrefix+"\t", container)
		if err != nil {
			return "", err
		}
		bodyLines = append(bodyLines, stmt2)
	}
	return text + strings.Join(bodyLines, "\n"), nil
}

func RenderArguments(arguments pyast.Arguments, container string) (string, error) {
	argTexts := []string{}
	for _, a := range arguments.Args() {
		text := a.Arg().Value()
		annotation := a.Annotation()
		if annotation == nil {
			return "", errors.New("TODO annotation expected for now...")
		}
		if annotation != nil {
			argText, err := renderStatement(pyastm.NewExprStatement(annotation), "", container)
			if err != nil {
				return "", nil
			}
			if argText == "" {
				return "", fmt.Errorf("Annotation with no text: v+%", annotation)
			}
			text = text + ": " + argText
		}

		argTexts = append(argTexts, text)
	}

	return "(" + strings.Join(argTexts, ", ") + ")", nil
}

// TODO Refactor together with the go version
func GeneratePySourceForDL(projectDir string, logger klog.Logger, dryRun bool) error {
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
