package parser

import (
	"fmt"
)

type ParseError struct {
	Message   string
	LineNum   int
	ColumnNum int
	Cause     error
}

func NewParseErrorAsError(message string, lineNum int, columnNum int, cause error) error {
	return NewParseError(message, lineNum, columnNum, cause)
}

func NewParseError(message string, lineNum int, columnNum int, cause error) *ParseError {
	return &ParseError{
		Message:   message,
		LineNum:   lineNum,
		ColumnNum: columnNum,
		Cause:     cause,
	}
}

func (e *ParseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%d:%d Parse error: %s; Cause: %s", e.LineNum, e.ColumnNum, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%d:%d Parse error: %s", e.LineNum, e.ColumnNum, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Cause
}
