package parser

import (
	"io"
	"os"

	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
)

type Parser interface {
	Parse(path, namespace string) (model.Design, *ParseError)
}

func NewParser() Parser {
	// TODO This got too complicated. Get an ANTLR expert to fix this.
	// return NewANTLRParser()

	return NewSinglePassParser()
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}
