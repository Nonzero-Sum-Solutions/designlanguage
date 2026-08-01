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
	currComponentLines          []string
	components                  []model.Component
}

func (c *spParserContext) newScanner() *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(c.document))
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

func (*singlePassParser) Parse(path, namespace string) (model.Design, error) {
	fileBytes, err := readFile(path)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("Error reading stream: %w", err)
	}

	if lineNum < 2 {
		return nil, errors.New("Empty document")
	}

	author, docComment := "", ""
	if prelude := pContext.docPrelude; prelude != nil {
		author = prelude.author
		docComment = prelude.comment
	}

	return model.NewDesign(author, docComment, namespace, pContext.components), nil
}

func parseLine(pContext *spParserContext, line string, lineNum int) error {
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
		if err := parseInComponent(line, pContext); err != nil {
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

	return fmt.Errorf("%d: Unexpected state. Line: [%s], Context: %v", lineNum, line, pContext)
}

func parseAfterComponentSpace(line string, lineNum int, pContext *spParserContext) error {
	if line == "" {
		return fmt.Errorf("%d: Excess space. line: [%s]", lineNum, line)
	}

	if IsComponentLine1(line) {
		if IsEnumLine(line) {
			enum1, err := ParseEnum(line)
			if err != nil {
				return err
			}
			pContext.components = append(pContext.components, enum1)
			pContext.lastLineType = EnumLine
			pContext.containerType = Start
			return nil
		}

		pContext.currComponentLines = []string{line}
		pContext.lastLineType = ComponentLine1
		pContext.containerType = Component
		return nil
	}

	return fmt.Errorf("%d: parseAfterComponentSpace: Expected component line 1 but got line: [%s]", lineNum, line)
}

func parseInComponent(line string, pContext *spParserContext) error {
	if line != "" {
		pContext.currComponentLines = append(pContext.currComponentLines, line)
		return nil
	}

	newComponent, err := parseComponent(pContext.currComponentLines)
	if err != nil {
		return err
	}

	pContext.components = append(pContext.components, newComponent)
	pContext.containerType = Start
	pContext.lastLineType = ComponentSpace
	pContext.currComponentLines = []string{}
	return nil
}

func parseAfterDocCommentSpace(line string, pContext *spParserContext, lineNum int) error {
	if IsComponentLine1(line) {
		if IsEnumLine(line) {
			enum1, err := ParseEnum(line)
			if err != nil {
				return err
			}
			pContext.components = append(pContext.components, enum1)
			pContext.lastLineType = ComponentSpace // EnumLine
			pContext.containerType = Start
			return nil
		}
		pContext.currComponentLines = []string{line}
		pContext.lastLineType = ComponentLine1
		pContext.containerType = Component
		return nil
	}

	return fmt.Errorf("%d: parseAfterDocCommentSpace: Expected component line 1 but got line: %s", lineNum, line)
}

func parseAfterDocComment(line string, pContext *spParserContext, lineNum int) error {
	if line == "" {
		pContext.lastLineType = DocCommentSpace
		return nil
	}

	return fmt.Errorf("%d: Expected empty line but got: %s", lineNum, line)
}

func parseAfterAuthor(line string, pContext *spParserContext) (shouldContinue bool) {
	if isComment(line) {
		pContext.docPrelude.comment = line[3:]
		pContext.lastLineType = DocComment
		return true
	}

	return false
}

func parseLine1(line string, pContext *spParserContext, lineNum int) error {
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
			enum1, err := ParseEnum(line)
			if err != nil {
				return err
			}
			pContext.components = append(pContext.components, enum1)
			pContext.lastLineType = EnumLine
			pContext.containerType = Start
			return nil
		}
		pContext.currComponentLines = []string{line}
		pContext.lastLineType = ComponentLine1
		pContext.containerType = Component
		return nil
	}

	return fmt.Errorf("%d: parseLine1: Expected component line 1 but got: [%s]", lineNum, line)
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

func parseComponent(lines []string) (model.Component, error) {
	if len(lines) < 1 {
		return nil, errors.New("Too few lines for a component")
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
		return nil, err
	}

	if len(lines) == 1 {
		return model.NewComponent(name, "", supertypeType)
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

	attrs, methods, err := parseAttributesAndMethods(fields)
	if err != nil {
		return nil, err
	}

	if len(attrs) == 0 && len(methods) == 0 {
		return model.NewComponent(name, comment, supertypeType)
	}

	if len(methods) > 0 {
		return model.NewObject(name, comment, supertypeType, attrs, methods)
	}

	return model.NewEntity(name, comment, supertypeType, attrs)
}

func parseAttributesAndMethods(lines []string) ([]model.Attribute, []model.Method, error) {
	methods := []model.Method{}
	attrs := []model.Attribute{}

	for _, line := range lines {
		if !strings.HasPrefix(line, "* ") {
			return nil, nil, errors.New("Component line doesn't start with '*'.")
		}

		content := line[2:]

		if strings.Contains(content, "(") {
			method, err := ParseMethod(content)
			if err != nil {
				return nil, nil, err
			}

			methods = append(methods, method)
			continue
		}

		attr, err := parseAttribute(content)
		if err != nil {
			return nil, nil, err
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

func ParseMethod(methodText string) (model.Method, error) {
	text := strings.TrimSpace(methodText)
	leftParensIdx1 := strings.Index(text, "(")
	if leftParensIdx1 < 1 {
		return nil, fmt.Errorf("Missing method left parens in text: %s", text)
	}

	beforeLP := text[:leftParensIdx1]

	if beforeLP[len(beforeLP)-1:] != " " {
		return nil, fmt.Errorf("Missing expected space before params. Text: [%s]", text)
	}

	name := beforeLP[:len(beforeLP)-1]

	rightParensIdx1 := strings.Index(text, ")")
	if rightParensIdx1 < leftParensIdx1+1 {
		return nil, fmt.Errorf("Mismatched parens in text: %s", text)
	}

	paramsText := text[leftParensIdx1+1 : rightParensIdx1]
	params, err := parseParams(paramsText)

	if rightParensIdx1 == len(text)-1 {
		return model.NewMethod(name, "", params, []model.Param{})
	}

	if len(text) < rightParensIdx1+3 {
		return nil, fmt.Errorf("Missing return vals in text: %s", text)
	}

	afterRightParens1Exp := text[rightParensIdx1+1:]

	if len(afterRightParens1Exp) < 4 || !strings.HasPrefix(afterRightParens1Exp, " -> ") {
		return nil, fmt.Errorf("Missing arrow: %s", afterRightParens1Exp)
	}

	afterArrowExp := afterRightParens1Exp[4:]
	afterNameTokens := strings.Split(afterArrowExp, " -- ")

	if len(afterNameTokens) < 1 {
		return nil, fmt.Errorf("Missing return expression in expression: %s", afterRightParens1Exp)
	}
	returnExp := afterNameTokens[0]
	comment := ""
	if len(afterNameTokens) == 1 {
		// No comment
	} else if len(afterNameTokens) > 2 {
		return nil, fmt.Errorf("Malformed return expression: %s", afterRightParens1Exp)
	} else {
		comment = afterNameTokens[1]
	}

	returnVals, err := parseParams(returnExp)
	if err != nil {
		return nil, err
	}

	return model.NewMethod(name, comment, params, returnVals)
}

// TODO Need to use this somewhere
func ParseEnum(text string) (model.Enum, error) {
	if !strings.Contains(text, " = {") {
		return nil, fmt.Errorf("No enum match: %s", text)
	}

	tokens := strings.Split(text, " = {")

	if len(tokens) != 2 {
		return nil, fmt.Errorf("Too few tokens: %s", text)
	}

	name := tokens[0]

	values := strings.Split(tokens[1][:len(tokens[1])-1], ", ")

	return model.NewEnum(name, values...)
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
