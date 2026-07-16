package pytranslator

import (
	"github.com/iancoleman/strcase"
	"github.com/mattmunz/appkit/maps"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/model/gen/pyast"
	pyastm "github.com/mattmunz/designlanguage/model/pyast"
)

var pyTypeByDLType = map[string]*pyType{
	"String": newPyType("", "", "str", false),
	"Int":    newPyType("", "", "int", false),
	"Error":  newPyType("", "", "Exception", false),
	// TODO Extract to config file
	"Logger": newPyType("", pyPackageNamesByAliasForImport["log"], "Logger", false),
}

func toPyFnCase(s string) string {
	return strcase.ToSnake(s)
}

func toPyArgCase(s string) string {
	return strcase.ToSnake(s)
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

// TODO Extract to config file
var pyAliasesByPackageNameForImport = map[string]string{
	"logging": "logger",
}

var pyPackageNamesByAliasForImport = maps.InvertMap(pyAliasesByPackageNameForImport)

func TranslateDesign(design model.Design) (pyast.AST, string, error) {
	ast1 := pyastm.NewAST()

	if design.Comment() != "" {
		ast1.Add(pyastm.NewCommentLine(design.Comment()))
	}

	// TODO IN future will need more sophistocated import mechanism
	ast1.Add(pyastm.NewImportFrom(1, pyastm.Id("typing"), []pyast.Alias{pyastm.NewAlias("Protocol")}))

	err := addDesignPy(ast1, design)
	if err != nil {
		return nil, "", err
	}
	return ast1, "", nil
}

func addDesignPy(ast1 pyast.AST, design model.Design) error {
	for _, e := range design.Enums() {
		addEnumPy(ast1, e)
	}

	for _, c := range design.BaseComponents() {
		AddBaseComponentPy(ast1, c)
	}

	for _, en := range design.Entities() {
		AddEntityPy(ast1, en)
	}

	for _, o := range design.Objects() {
		AddObjectPy(ast1, o)
	}

	return nil
}

func getPyType(type1 model.Type) *pyType {
	if pyType, exists := pyTypeByDLType[type1.Name()]; exists {
		return newPyType(pyType.op, pyType.packageName, pyType.name, type1.IsArray())
	}
	return newPyType("", "", type1.Name(), type1.IsArray())
}

func AddObjectPy(ast1 pyast.AST, obj model.Object) error {
	body := []pyast.Statement{}
	if obj.Comment() != "" {
		comment := pyastm.NewCommentLine(obj.Comment())
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(comment)))
	}
	if len(obj.Attributes())+len(obj.Methods()) == 0 {
		// Should be impossible
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

func NewMethodDef(m model.Method) (pyast.FunctionDef, error) {
	var returnsStmt pyast.Expression
	var err error

	returnVals := m.ReturnVals()

	raisesException := false
	elts := []pyast.Expression{}
	for _, v := range returnVals {
		if v.Type().Name() == "Error" {
			raisesException = true
			continue
		}

		var elt pyast.Expression
		pyType := getPyType(v.Type())
		if !pyType.isArray {
			elt = pyastm.NewName(pyType.name, pyast.Load)
		} else {
			elt = pyastm.NewName("List<"+pyType.name+">", pyast.Load)
		}
		elts = append(elts, elt)
	}

	body := []pyast.Statement{}
	if m.Comment() != "" || raisesException {
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

func AddEntityPy(ast1 pyast.AST, entity model.Entity) {
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
	body := []pyast.Statement{}
	if a.Comment() != "" {
		comment := pyastm.NewCommentLine(a.Comment())
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(comment)))
	}
	body = append(body, pyastm.NewExprStatement(pyastm.NewEllipsis()))

	pyType := getPyType(a.Type())
	args := pyastm.NewArguments([]pyast.Arg{})
	return pyastm.NewFunctionDef(toPyFnCase(a.Name()), args, body, toPYASTTypeExpression(pyType))
}

func toPYASTTypeExpression(pyType *pyType) pyast.Expression {
	if !pyType.isArray {
		return pyastm.NewName(pyType.name, pyast.Load)
	}
	return pyastm.NewName("List<"+pyType.name+">", pyast.Load)
}

func addEnumPy(ast1 pyast.AST, e model.Enum) {
	panic("TODO NYI aep")
}

func AddBaseComponentPy(ast1 pyast.AST, component model.Component) {
	body := []pyast.Statement{}
	if component.Comment() != "" {
		body = append(body, pyastm.NewExprStatement(pyastm.NewDocString(pyastm.NewCommentLine(component.Comment()))))
	}
	body = append(body, pyastm.NewExprStatement(pyastm.NewEllipsis()))

	ast1.Add(pyastm.NewClassDef(pyastm.Id(component.Name()), getClassDefBases(component), nil, nil, nil, body, nil))
}
