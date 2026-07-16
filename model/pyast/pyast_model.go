package pyast

import (
	"errors"

	"github.com/mattmunz/designlanguage/model/gen/pyast"
	pyastg "github.com/mattmunz/designlanguage/model/gen/pyast"
)

type node struct{}

func (n *node) Position() pyastg.Position {
	panic("unimplemented")
}

func NewNode() pyastg.Node {
	return &node{}
}

type identifier struct {
	value string
}

func (i *identifier) Value() string {
	return i.value
}

func Id(name string) pyastg.Identifier {
	return &identifier{name}
}

type expression struct {
	body           []pyastg.Statement
	expressionType string
}

// ExpressionType implements [pyast.Expression].
func (e *expression) ExpressionType() string {
	return e.expressionType
}

func (e *expression) Body() []pyastg.Statement {
	return e.body
}

func (e *expression) Position() pyastg.Position {
	panic("unimplemented 12")
}

func NewExpression(body []pyastg.Statement, expressionType string) pyastg.Expression {
	return &expression{body, expressionType}
}

type exprStatement struct {
	pyast.Statement
	value pyastg.Expression
}

func (e *exprStatement) ExprValue() pyastg.Expression {
	return e.value
}

func NewExprStatement(value pyastg.Expression) pyast.ExprStatement {
	return &exprStatement{NewStatement("ExprStatement"), value}
}

type alias struct {
	pyast.Node
	name, asName pyastg.Identifier
}

func (a *alias) AsName() pyastg.Identifier {
	return a.asName
}

func (a *alias) Name() pyastg.Identifier {
	return a.name
}

func NewAlias(name string) pyast.Alias {
	return &alias{NewNode(), Id(name), nil}
}

type importFrom struct {
	pyastg.Statement
	level  int
	module pyastg.Identifier
	names  []pyastg.Alias
}

func (i *importFrom) Level() int {
	return i.level
}

func (i *importFrom) Module() pyastg.Identifier {
	return i.module
}

func (i *importFrom) Names() []pyastg.Alias {
	return i.names
}

func NewImportFrom(level int, module pyastg.Identifier, names []pyastg.Alias) pyastg.ImportFrom {
	return &importFrom{NewStatement("ImportFrom"), level, module, names}
}

type tuple struct {
	pyast.Expression
	ctx  pyastg.ExpressionContext
	elts []pyastg.Expression
}

func (t *tuple) Ctx() pyastg.ExpressionContext {
	return t.ctx
}

func (t *tuple) Elts() []pyastg.Expression {
	return t.elts
}

func NewTuple(elts []pyastg.Expression, ctx pyastg.ExpressionContext) (pyast.Tuple, error) {
	if len(elts) < 1 {
		return nil, errors.New("At least 1 elt expected")
	}
	return &tuple{NewExpression([]pyast.Statement{}, "Tuple"), ctx, elts}, nil
}

type ast struct {
	nodes []pyastg.Node
}

func (a *ast) Add(node pyastg.Node) {
	a.nodes = append(a.nodes, node)
}

func (a *ast) Nodes() []pyastg.Node {
	return a.nodes
}

func (a *ast) Position() pyastg.Position {
	panic("unimplemented 3")
}

func NewAST() pyastg.AST {
	return &ast{}
}

type docString struct {
	pyast.Expression
}

func NewDocString(lines ...pyast.CommentLine) pyastg.DocString {
	return &docString{NewExpression(toStatements(lines), "DocString")}
}

func toStatements(lines []pyastg.CommentLine) []pyastg.Statement {
	statements := make([]pyastg.Statement, len(lines))
	for i, l := range lines {
		statements[i] = l
	}
	return statements
}

type returnStmt struct {
	pyastg.Statement
	value pyast.Expression
}

func (r *returnStmt) ReturnValue() pyastg.Expression {
	return r.value
}

// AST: returns=Name(id='foo', ctx=Load()))]
func NewReturn(value pyast.Expression) pyastg.Return {
	return &returnStmt{NewStatement("Return"), value}
}

type arg struct {
	pyast.Node
	arg        pyastg.Identifier
	annotation pyastg.Expression
}

func (a *arg) Annotation() pyastg.Expression {
	return a.annotation
}

func (a *arg) Arg() pyastg.Identifier {
	return a.arg
}

func NewArg(arg1 pyastg.Identifier, annotation pyastg.Expression) pyastg.Arg {
	return &arg{NewNode(), arg1, annotation}
}

type arguments struct {
	args []pyastg.Arg
}

func (a *arguments) Args() []pyastg.Arg {
	return a.args
}

// Defaults implements [pyast.Arguments].
func (a *arguments) Defaults() []pyastg.Expression {
	panic("unimplemented")
}

// KwDefaults implements [pyast.Arguments].
func (a *arguments) KwDefaults() []pyastg.Expression {
	panic("unimplemented")
}

// Kwarg implements [pyast.Arguments].
func (a *arguments) Kwarg() pyastg.Arg {
	panic("unimplemented")
}

// Kwonlyargs implements [pyast.Arguments].
func (a *arguments) Kwonlyargs() []pyastg.Arg {
	panic("unimplemented")
}

// Position implements [pyast.Arguments].
func (a *arguments) Position() pyastg.Position {
	panic("unimplemented")
}

// Vararg implements [pyast.Arguments].
func (a *arguments) Vararg() pyastg.Arg {
	panic("unimplemented")
}

func NewArguments(args []pyastg.Arg) pyast.Arguments {
	return &arguments{args}
}

type functionDef struct {
	pyastg.Statement
	name          pyast.Identifier
	args          pyast.Arguments
	body          []pyast.Statement
	decoratorList []pyast.Expression
	returns       pyast.Expression
}

// Args implements [pyast.FunctionDef].
func (f functionDef) Args() pyastg.Arguments {
	return f.args
}

// Body implements [pyast.FunctionDef].
func (f functionDef) Body() []pyastg.Statement {
	return f.body
}

// DecoratorList implements [pyast.FunctionDef].
func (f functionDef) DecoratorList() []pyastg.Expression {
	return f.decoratorList
}

// Name implements [pyast.FunctionDef].
func (f functionDef) Name() pyastg.Identifier {
	return f.name
}

// Returns implements [pyast.FunctionDef].
func (f functionDef) Returns() pyastg.Expression {
	return f.returns
}

func NewFunctionDef(name string, args pyast.Arguments, body []pyastg.Statement, returns pyastg.Expression) pyastg.FunctionDef {
	return functionDef{NewStatement("FunctionDef"), Id(name), args, body, nil, returns}
}

type classDef struct {
	pyastg.Statement
	name          pyastg.Identifier
	bases         []pyastg.Expression
	keywords      []pyastg.KeywordNode
	starargs      pyastg.Expression
	kwargs        pyastg.Expression
	body          []pyastg.Statement
	decoratorList []pyastg.Expression
}

func (c *classDef) Position() pyastg.Position {
	panic("unimplemented 4")
}

func (c *classDef) Bases() []pyastg.Expression {
	return c.bases
}

func (c *classDef) Body() []pyastg.Statement {
	return c.body
}

func (c *classDef) DecoratorList() []pyastg.Expression {
	panic("unimplemented 7")
}

func (c *classDef) Keywords() []pyastg.KeywordNode {
	panic("unimplemented 8")
}

func (c *classDef) Kwargs() pyastg.Expression {
	panic("unimplemented 9")
}

func (c *classDef) Name() pyastg.Identifier {
	return c.name
}

func (c *classDef) Starargs() pyastg.Expression {
	panic("unimplemented 11")
}

func NewClassDef(name pyastg.Identifier, bases []pyastg.Expression, keywords []pyastg.KeywordNode, starargs pyastg.Expression,
	kwargs pyastg.Expression, body []pyastg.Statement, decoratorList []pyastg.Expression) pyastg.ClassDef {
	return &classDef{NewStatement("ClassDef"), name, bases, keywords, starargs, kwargs, body, decoratorList}
}

type name struct {
	pyastg.Expression
	id  pyastg.Identifier
	ctx pyastg.ExpressionContext
}

func (n *name) Ctx() pyastg.ExpressionContext {
	return n.ctx
}

func (n *name) Id() pyastg.Identifier {
	return n.id
}

func NewName(name1 string, ctx pyastg.ExpressionContext) pyastg.Name {
	return &name{NewExpression([]pyastg.Statement{}, "Name"), Id(name1), ctx}
}

type ellipsis struct {
	pyastg.Expression
	symbol string
}

func (e *ellipsis) Symbol() string {
	return e.symbol
}

func NewEllipsis() pyastg.Ellipsis {
	return &ellipsis{NewExpression([]pyastg.Statement{}, "Ellipsis"), "..."}
}

type statement struct {
	pyastg.Node
	statementType string
}

func (s *statement) StatementType() string {
	return s.statementType
}

func NewStatement(statementType string) pyastg.Statement {
	return &statement{NewNode(), statementType}
}

type commentLine struct {
	pyast.Statement
	text string
}

func (c *commentLine) Text() string {
	return c.text
}

func NewCommentLine(text string) pyastg.CommentLine {
	return &commentLine{NewStatement("CommentLine"), text}
}
