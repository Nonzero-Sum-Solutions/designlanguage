package parser

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
)

type ParserState string

const (
	Start           ParserState = ""
	Author          ParserState = "author"
	DocComment      ParserState = "docComment"
	DocCommentSpace ParserState = "docCommentSpace"
	Component       ParserState = "component"
	ComponentSpace  ParserState = "componentSpace"
	ComponentLine1  ParserState = "componentLine1"
	EnumLine        ParserState = "enumLine"
)

var nameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z]*$`)

type spParserContext struct {
	document                    string
	lastLineType, containerType ParserState
	docPrelude                  *prelude
	// TODO Is this allowed to have an empty line at the end?
	currComponentLines []string
	components         []model.Component
}

func (c *spParserContext) newScanner() *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(c.document))
}

func (c *spParserContext) SetCurrComponentLines(lines []string) {
	c.currComponentLines = lines
}

func (c *spParserContext) AppendCurrComponentLine(line string) {
	c.currComponentLines = append(c.currComponentLines, line)
	if strings.TrimSpace(line) == "" {
		panic(fmt.Sprintf("TODO Unexpected empty line in component. lines: [%+v]", c.currComponentLines))
	}
}

func (c *spParserContext) ClearCurrComponentLines() {
	c.currComponentLines = []string{}
}

func newSPParserContext(document string) *spParserContext {
	return &spParserContext{document, Start, Start, nil, []string{}, []model.Component{}}
}

type singlePassParser struct {
	context *spParserContext
}

type prelude struct {
	author, comment string
}

func NewSinglePassParser() Parser {
	return &singlePassParser{}
}

func (*singlePassParser) Parse(path, namespace string) (model.Design, *ParseError) {
	fileBytes, err := readFile(path)
	if err != nil {
		return nil, NewParseError("Couldn't read file", 0, 0, err)
	}

	pContext := newSPParserContext(string(fileBytes) + "\n")
	scanner := pContext.newScanner()
	lineNum := 0
	for lineNum = 1; scanner.Scan(); lineNum++ {
		line := scanner.Text()
		if err := parseLine(pContext, line, lineNum); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, NewParseError("Error reading stream", lineNum, 0, err)
	}

	if lineNum < 2 {
		return nil, NewParseError("Empty document", lineNum, 0, nil)
	}

	lineCount := strings.Count(pContext.document, "\n")
	if lineNum < lineCount {
		return nil, NewParseError(fmt.Sprintf("Too few lines in document processed. Expected at least %d", lineCount), lineNum, 0, nil)
	}

	author, docComment := "", ""
	if prelude := pContext.docPrelude; prelude != nil {
		author = prelude.author
		docComment = prelude.comment
	}

	return model.NewDesign(author, docComment, namespace, pContext.components), nil
}

func parseLine(pContext *spParserContext, line string, lineNum int) *ParseError {
	if pContext.lastLineType == Start {
		if err := parseLine1(line, pContext, lineNum); err != nil {
			return err
		}

		return nil
	}

	if pContext.lastLineType == Author && parseAfterAuthor(line, pContext) {
		return nil
	}

	if pContext.lastLineType == DocComment {
		if err := parseAfterDocComment(line, pContext, lineNum); err != nil {
			return err
		}
		return nil
	}

	if pContext.lastLineType == DocCommentSpace {
		if err := parseAfterDocCommentSpace(line, pContext, lineNum); err != nil {
			return err
		}
		return nil
	}

	if pContext.lastLineType == EnumLine && line == "" {
		pContext.lastLineType = ComponentSpace
		return nil
	}

	if pContext.containerType == Component {
		if err := parseInComponent(line, pContext, lineNum); err != nil {
			return err
		}
		return nil
	}

	if pContext.lastLineType == ComponentSpace {
		if err := parseAfterComponentSpace(line, lineNum, pContext); err != nil {
			return err
		}
		return nil
	}

	return NewParseError(fmt.Sprintf("Unexpected state. Context: %v", pContext), lineNum, 0, nil)
}

func parseAfterComponentSpace(line string, lineNum int, pContext *spParserContext) *ParseError {
	if line == "" {
		return NewParseError(fmt.Sprintf("Excess space. line: [%s]", line), lineNum, 0, nil)
	}

	if IsComponentLine1(line) {
		if IsEnumLine(line) {
			enum1, err := ParseEnum(line, lineNum)
			if err != nil {
				return err
			}
			pContext.components = append(pContext.components, enum1)
			pContext.lastLineType = EnumLine
			pContext.containerType = Start
			return nil
		}

		pContext.ClearCurrComponentLines()
		pContext.AppendCurrComponentLine(line)
		pContext.lastLineType = ComponentLine1
		pContext.containerType = Component
		return nil
	}

	return NewParseError(fmt.Sprintf("Expected component line 1 but got line: [%s]", line), lineNum, 0, nil)
}

func parseInComponent(line string, pContext *spParserContext, lineNum int) *ParseError {
	if line != "" {
		pContext.AppendCurrComponentLine(line)
		return nil
	}

	newComponent, err := parseComponent(pContext.currComponentLines, lineNum)
	if err != nil {
		return err
	}

	pContext.components = append(pContext.components, newComponent)
	pContext.containerType = Start
	pContext.lastLineType = ComponentSpace
	pContext.ClearCurrComponentLines()
	return nil
}

func parseAfterDocCommentSpace(line string, pContext *spParserContext, lineNum int) *ParseError {
	if IsComponentLine1(line) {
		if IsEnumLine(line) {
			enum1, err := ParseEnum(line, lineNum)
			if err != nil {
				return err
			}
			pContext.components = append(pContext.components, enum1)
			pContext.lastLineType = ComponentSpace // EnumLine
			pContext.containerType = Start
			return nil
		}
		pContext.ClearCurrComponentLines()
		pContext.AppendCurrComponentLine(line)
		pContext.lastLineType = ComponentLine1
		pContext.containerType = Component
		return nil
	}

	return NewParseError(fmt.Sprintf("Expected component line 1 but got line: [%s]", line), lineNum, 0, nil)
}

func parseAfterDocComment(line string, pContext *spParserContext, lineNum int) *ParseError {
	if line == "" {
		pContext.lastLineType = DocCommentSpace
		return nil
	}

	return NewParseError(fmt.Sprintf("Expected empty line but got: [%s]", line), lineNum, 0, nil)
}

func parseAfterAuthor(line string, pContext *spParserContext) (shouldContinue bool) {
	if isComment(line) {
		pContext.docPrelude.comment = line[3:]
		pContext.lastLineType = DocComment
		return true
	}

	return false
}

func parseLine1(line string, pContext *spParserContext, lineNum int) *ParseError {
	// Either this is author, line comment, or first component, or error
	if isComment(line) {
		if strings.HasPrefix(line, "-- Author: ") {
			if pContext.docPrelude == nil {
				pContext.docPrelude = &prelude{}
			}
			pContext.docPrelude.author = line[11:]
			pContext.lastLineType = Author
			return nil
		}

		if pContext.docPrelude == nil {
			pContext.docPrelude = &prelude{}
		}
		pContext.docPrelude.comment = line[3:]
		pContext.lastLineType = DocComment
		return nil
	}

	if IsComponentLine1(line) {
		if IsEnumLine(line) {
			enum1, err := ParseEnum(line, lineNum)
			if err != nil {
				return err
			}
			pContext.components = append(pContext.components, enum1)
			pContext.lastLineType = EnumLine
			pContext.containerType = Start
			return nil
		}
		pContext.ClearCurrComponentLines()
		pContext.AppendCurrComponentLine(line)
		pContext.lastLineType = ComponentLine1
		pContext.containerType = Component
		return nil
	}

	return NewParseError(fmt.Sprintf("Expected component line 1 but got: [%s]", line), lineNum, 0, nil)
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "-- ")
}

func IsComponentLine1(line string) bool {
	if IsEnumLine(line) {
		return true
	}

	tokens := strings.Split(line, " :: ")
	if len(tokens) < 1 || len(tokens) > 2 {
		return false
	}

	name := tokens[0]
	if len(tokens) < 2 {
		return nameRegex.MatchString(name)
	}

	return nameRegex.MatchString(name) && nameRegex.MatchString(tokens[1])
}

func IsEnumLine(line string) bool {
	return strings.Contains(line, " = {") && strings.HasSuffix(line, "}")
}

// TODO Too complicatedd. Refactor
func parseComponent(lines []string, lineNum int) (model.Component, *ParseError) {
	if len(lines) < 1 {
		return nil, NewParseError("Too few lines for a component", lineNum, 0, nil)
	}

	if strings.TrimSpace(lines[len(lines)-1]) == "" {
		return nil, NewParseError("Empty line not permitted at the end of the component", lineNum, 0, nil)
	}

	line1 := lines[0]
	name := ""
	supertype := ""
	line1Tokens := strings.Split(line1, " :: ")
	name = line1Tokens[0]
	if len(line1Tokens) > 1 {
		supertype = line1Tokens[1]
	}

	supertypeType, err := parseOptionalType(supertype)
	if err != nil {
		return nil, NewParseError("Couldn't parse supertype", lineNum, 0, err)
	}

	if len(lines) == 1 {
		component, err := model.NewComponent(name, "", supertypeType)
		if err != nil {
			return nil, NewParseError(fmt.Sprintf("Couldn't create component: [%s]", name), lineNum, 0, err)
		}
		return component, nil
	}

	comment := ""
	fields := []string{}

	line2 := lines[1]
	if isComment(line2) {
		comment = line2[3:]
	} else {
		fields = append(fields, line2)
	}

	if len(lines) > 2 {
		fields = append(fields, lines[2:]...)
	}

	// TODO On behalf of component. But are we including trailing blank lines?
	attrs, methods, err2 := parseAttributesAndMethods(fields, lineNum)
	if err2 != nil {
		return nil, err2
	}

	if len(attrs) == 0 && len(methods) == 0 {
		component, err := model.NewComponent(name, comment, supertypeType)
		if err != nil {
			return nil, NewParseError(fmt.Sprintf("Couldn't create component: [%s]", name), lineNum, 0, err)
		}
		return component, nil
	}

	if len(methods) > 0 {
		object, err := model.NewObject(name, comment, supertypeType, attrs, methods)
		if err != nil {
			return nil, NewParseError(fmt.Sprintf("Couldn't create object: [%s]", name), lineNum, 0, err)
		}
		return object, nil
	}

	entity, err := model.NewEntity(name, comment, supertypeType, attrs)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("Couldn't create entity: [%s]", name), lineNum, 0, err)
	}
	return entity, nil
}

func parseAttributesAndMethods(lines []string, lineNum int) ([]model.Attribute, []model.Method, *ParseError) {
	methods := []model.Method{}
	attrs := []model.Attribute{}

	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		return nil, nil, NewParseError("Unexpected Empty line not at the end of the component. Not sure if this is a real error or not", lineNum, 0, nil)
	}

	for i, line := range lines {
		curLineNum := lineNum + i
		if !strings.HasPrefix(line, "* ") {
			// TODO This is called inappropriately.
			msg := fmt.Sprintf("Component line doesn't start with '*'. Line [%d] out of [%d]", i, len(lines))
			return nil, nil, NewParseError(msg, curLineNum, 0, nil)
		}

		content := line[2:]

		if strings.Contains(content, "(") {
			method, err := ParseMethod(content, curLineNum)
			if err != nil {
				return nil, nil, err
			}

			methods = append(methods, method)
			continue
		}

		attr, err := parseAttribute(content)
		if err != nil {
			return nil, nil, NewParseError("Couldn't parse attribute", curLineNum, 0, err)
		}
		attrs = append(attrs, attr)
	}

	return attrs, methods, nil
}

func parseAttribute(text string) (model.Attribute, error) {
	tokens := strings.Split(text, " ")
	if len(tokens) < 2 {
		return nil, errors.New("Expected at least 2 tokens")
	}
	name := tokens[0]
	typeExp := tokens[1]
	typ, err := parseType(typeExp)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 2 {
		return model.NewAttribute(name, "", typ.Name(), typ.IsArray())
	}

	comment, err := parseEndOfLineComment(tokens[2:], text)
	if err != nil {
		return nil, err
	}

	return model.NewAttribute(name, comment, typ.Name(), typ.IsArray())
}

func ParseMethod(methodText string, lineNum int) (model.Method, *ParseError) {
	text := strings.TrimSpace(methodText)
	leftParensIdx1 := strings.Index(text, "(")
	if leftParensIdx1 < 1 {
		return nil, NewParseError(fmt.Sprintf("Missing method left parens in text: [%s]", text), lineNum, 0, nil)
	}

	beforeLP := text[:leftParensIdx1]

	if beforeLP[len(beforeLP)-1:] != " " {
		return nil, NewParseError(fmt.Sprintf("Missing expected space before params. Text: [%s]", text), lineNum, 0, nil)
	}

	name := beforeLP[:len(beforeLP)-1]

	rightParensIdx1 := strings.Index(text, ")")
	if rightParensIdx1 < leftParensIdx1+1 {
		return nil, NewParseError(fmt.Sprintf("Mismatched parens in text: [%s]", text), lineNum, 0, nil)
	}

	paramsText := text[leftParensIdx1+1 : rightParensIdx1]
	params, err := parseParams(paramsText)

	if rightParensIdx1 == len(text)-1 {
		method, err := model.NewMethod(name, "", params, []model.Param{})
		if err != nil {
			return nil, NewParseError(fmt.Sprintf("Couldn't create method: [%s]", name), lineNum, 0, err)
		}
		return method, nil
	}

	if len(text) < rightParensIdx1+3 {
		return nil, NewParseError(fmt.Sprintf("Missing return vals in text: [%s]", text), lineNum, 0, nil)
	}

	afterRightParens1Exp := text[rightParensIdx1+1:]

	if len(afterRightParens1Exp) < 4 || !strings.HasPrefix(afterRightParens1Exp, " -> ") {
		return nil, NewParseError(fmt.Sprintf("Missing arrow: [%s]", afterRightParens1Exp), lineNum, 0, nil)
	}

	afterArrowExp := afterRightParens1Exp[4:]
	afterNameTokens := strings.Split(afterArrowExp, " -- ")

	if len(afterNameTokens) < 1 {
		return nil, NewParseError(fmt.Sprintf("Missing return expression in expression: [%s]", afterRightParens1Exp), lineNum, 0, nil)
	}
	returnExp := afterNameTokens[0]
	comment := ""
	if len(afterNameTokens) == 1 {
		// No comment
	} else if len(afterNameTokens) > 2 {
		return nil, NewParseError(fmt.Sprintf("Malformed return expression: [%s]", afterRightParens1Exp), lineNum, 0, nil)
	} else {
		comment = afterNameTokens[1]
	}

	returnVals, err := parseParams(returnExp)
	if err != nil {
		return nil, NewParseError("Couldn't parse return values", lineNum, 0, err)
	}

	method, err := model.NewMethod(name, comment, params, returnVals)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("Couldn't create method: [%s]", name), lineNum, 0, err)
	}
	return method, nil
}

// TODO Need to use this somewhere
func ParseEnum(text string, lineNum int) (model.Enum, *ParseError) {
	if !strings.Contains(text, " = {") {
		return nil, NewParseError(fmt.Sprintf("No enum match: [%s]", text), lineNum, 0, nil)
	}

	tokens := strings.Split(text, " = {")

	if len(tokens) != 2 {
		return nil, NewParseError(fmt.Sprintf("Too few tokens: [%s]", text), lineNum, 0, nil)
	}

	name := tokens[0]

	values := strings.Split(tokens[1][:len(tokens[1])-1], ", ")

	enum, err := model.NewEnum(name, values...)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("Couldn't create enum: [%s]", text), lineNum, 0, err)
	}
	return enum, nil
}

// Input may either be in one of the following forms.
// ()
// ""
// (A B, C D, E F)
// A B, C D, E F
func parseParams(text string) ([]model.Param, error) {
	text2 := strings.ReplaceAll(text, "(", "")
	text3 := strings.ReplaceAll(text2, ")", "")
	tokens := strings.Split(text3, ", ")

	params := []model.Param{}

	for _, paramExp := range tokens {
		if paramExp == "" {
			continue
		}

		param, err := parseParam(paramExp)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse paramExp: %s, Err: %w", paramExp, err)
		}

		params = append(params, param)
	}

	return params, nil
}

func parseParam(paramExp string) (model.Param, error) {
	tokens := strings.Split(paramExp, " ")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("Wrong number of param elements: %s", paramExp)
	}
	t, err := parseType(tokens[1])
	return model.NewParam(tokens[0], t), err
}

func parseEndOfLineComment(tokens []string, text string) (string, error) {
	if len(tokens) < 2 {
		return "", fmt.Errorf("Expected more tokens for comment on line: %s", text)
	}

	if tokens[0] != "--" {
		return "", fmt.Errorf("Expected comment prefix: tokens: %s, text: %s", tokens, text)
	}

	return strings.Join(tokens[1:], " "), nil
}

func parseOptionalType(typeExp string) (model.Type, error) {
	if typeExp == "" {
		return nil, nil
	}

	return parseType(typeExp)
}

func parseType(typeExp string) (model.Type, error) {
	isArray := false
	name := typeExp
	if strings.HasPrefix(typeExp, "[]") {
		isArray = true
		name = typeExp[2:]
	}

	return model.NewType(name, isArray)
}
