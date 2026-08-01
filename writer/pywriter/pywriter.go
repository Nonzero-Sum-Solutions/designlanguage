package pywriter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nonzero-Sum-Solutions/designlanguage/model/gen/pyast"
	pyastm "github.com/Nonzero-Sum-Solutions/designlanguage/model/pyast"
)

func RenderPy(ast1 pyast.AST, linePrefix, container string) (string, error) {
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

// TODO Handle "asname" when needed
func renderImportFrom(i pyast.ImportFrom) string {
	nameTexts := []string{}
	for _, n := range i.Names() {
		nameTexts = append(nameTexts, n.Name().Value())
	}

	return "from " + i.Module().Value() + " import " + strings.Join(nameTexts, ", ")
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
	lines = append(lines, linePrefix+"\"\"\"\n")
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
