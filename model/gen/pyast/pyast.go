// A simple, workable AST for Python source code.
package pyast

type ExpressionContext int

const (
	Load ExpressionContext = iota
	Store
	Del
	AugLoad
	AugStore
	Param
)

type BoolOp int

const (
	And BoolOp = iota
	Or
)

type Op int

const (
	Add Op = iota
	Sub
	Mult
	Div
	Modulo
	Pow
	LShift
	RShift
	BitOr
	BitXor
	BitAnd
	FloorDiv
)

type UnaryOp int

const (
	Invert UnaryOp = iota
	Not
	UAdd
	USub
)

type BinaryOp int

const (
	Eq BinaryOp = iota
	NotEq
	Lt
	LtE
	Gt
	GtE
	Is
	IsNot
	In
	NotIn
)

type Keyword int

const (
	False Keyword = iota
	None
	True
	_and
	_as
	_assert
	_async
	_await
	_break
	_class
	_continue
	_def
	_del
	_elif
	_else
	_except
	_finally
	_for
	_from
	_global
	_if
	_import
	_in
	_is
	_lambda
	_nonlocal
	_not
	_or
	_pass
	_raise
	_return
	_try
	_while
	_with
	_yield
)

type Object interface{}

type Singleton interface{}

type Bytes interface {
	Node
}

type DocString interface {
	Expression
}

type Pass interface {
	Statement
}

type Break interface {
	Statement
}

type Continue interface {
	Statement
}

type SliceBase interface {
	Node
}

type Identifier interface {
	Value() string
}

type Position interface {
	Line() int
	ColOffset() int
}

type Node interface {
	Object

	Position() Position
}

type Module interface {
	Node

	Body() []Statement
}

type Expression interface {
	Node

	Body() []Statement
	ExpressionType() string
}

type Suite interface {
	Node

	Body() []Statement
}

type Statement interface {
	Node

	StatementType() string
}

type CommentLine interface {
	Statement

	Text() string
}

type FunctionDef interface {
	Statement

	Name() Identifier
	Args() Arguments
	Body() []Statement
	DecoratorList() []Expression
	// The return type. The return expression in the body is somethign ddifferent.
	Returns() Expression
}

type ClassDef interface {
	Statement

	Name() Identifier
	Bases() []Expression
	Keywords() []KeywordNode
	Starargs() Expression
	Kwargs() Expression
	Body() []Statement
	DecoratorList() []Expression
}

type Return interface {
	Statement

	ReturnValue() Expression
}

type Delete interface {
	Statement

	Targets() []Expression
}

type Assign interface {
	Statement

	Targets() []Expression
	Value() Expression
}

type AugAssign interface {
	Statement

	Target() Expression
	Op() Op
	Value() Expression
}

type For interface {
	Statement

	Target() Expression
	Iter() Expression
	Body() []Statement
	Orelse() []Statement
}

type While interface {
	Statement

	Test() Expression
	Body() []Statement
	Orelse() []Statement
}

type If interface {
	Statement

	Test() Expression
	Body() []Statement
	Orelse() []Statement
}

type With interface {
	Statement

	Items() []WithItem
	Body() []Statement
}

type Raise interface {
	Statement

	Exc() Expression
	Cause() Expression
}

type Try interface {
	Statement

	Body() []Statement
	Handlers() []ExceptHandler
	Orelse() []Statement
	Finalbody() []Statement
}

type Assert interface {
	Statement

	Test() Expression
	Msg() Expression
}

type Import interface {
	Statement

	Names() []Alias
}

type ImportFrom interface {
	Statement

	Module() Identifier
	Names() []Alias
	Level() int
}

type Global interface {
	Statement

	GlobalNames() []Identifier
}

type Nonlocal interface {
	Statement

	NonlocalNames() []Identifier
}

type ExprStatement interface {
	Statement

	ExprValue() Expression
}

type BoolOpExpr interface {
	Expression

	Op() BoolOp
	Values() []Expression
}

type BinOp interface {
	Expression

	Left() Expression
	Op() Op
	Right() Expression
}

type UnaryOpExp interface {
	Expression

	Op() UnaryOp
	Operand() Expression
}

type Lambda interface {
	Expression

	Args() Arguments
}

type IfExp interface {
	Expression

	Test() Expression
	Orelse() Expression
}

type Dict interface {
	Expression

	Keys() []Expression
	Values() []Expression
}

type Set interface {
	Expression

	Elts() []Expression
}

type ListComp interface {
	Expression

	Elt() Expression
	Generators() []Comprehension
}

type SetComp interface {
	Expression

	Elt() Expression
	Generators() []Comprehension
}

type DictComp interface {
	Expression

	Key() Expression
	Value() Expression
	Generators() []Comprehension
}

type GeneratorExp interface {
	Expression

	Elt() Expression
	Generators() []Comprehension
}

type Yield interface {
	Expression

	YieldValue() Expression
}

type YieldFrom interface {
	Expression

	YieldFromValue() Expression
}

type Compare interface {
	Expression

	Left() Expression
	Ops() []BinaryOp
	Comparators() []Expression
}

type Call interface {
	Expression

	Func() Expression
	Args() []Expression
	Keywords() []KeywordNode
	Starargs() Expression
	Kwargs() Expression
}

type Num interface {
	Expression

	N() Object
}

type Str interface {
	Expression

	S() string
}

type BytesExp interface {
	Expression

	S() Bytes
}

type NameConstant interface {
	Expression

	Value() Singleton
}

type Ellipsis interface {
	Expression

	Symbol() string
}

type Attribute interface {
	Expression

	Value() Expression
	Attr() Identifier
	Ctx() ExpressionContext
}

type Subscript interface {
	Expression

	Value() Expression
	Slice() Slice
	Ctx() ExpressionContext
}

type Starred interface {
	Expression

	Value() Expression
	Ctx() ExpressionContext
}

type Name interface {
	Expression

	Id() Identifier
	Ctx() ExpressionContext
}

type List interface {
	Expression

	Elts() []Expression
	Ctx() ExpressionContext
}

type Tuple interface {
	Expression

	Elts() []Expression
	Ctx() ExpressionContext
}

type Slice interface {
	SliceBase

	Lower() Expression
	Upper() Expression
	Step() Expression
}

type ExtSlice interface {
	SliceBase

	Dims() []Slice
}

type Index interface {
	SliceBase

	Value() Expression
}

type Comprehension interface {
	Expression

	Target() Expression
	Iter() Expression
	Ifs() []Expression
}

type ExceptHandler interface {
	Node

	ExpressionType() Expression
	Name() Identifier
	Body() []Statement
}

type Arguments interface {
	Node

	Args() []Arg
	Vararg() Arg
	Kwonlyargs() []Arg
	KwDefaults() []Expression
	Kwarg() Arg
	Defaults() []Expression
}

type Arg interface {
	Node

	Arg() Identifier
	Annotation() Expression
}

type KeywordNode interface {
	Node

	Arg() Identifier
	Value() Expression
}

type Alias interface {
	Node

	Name() Identifier
	AsName() Identifier
}

type WithItem interface {
	Node

	ContextExpression() Expression
	OptionalVars() Expression
}

type AST interface {
	Node

	Nodes() []Node
	Add(node Node)
}
