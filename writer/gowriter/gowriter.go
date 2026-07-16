package gowriter

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/dave/jennifer/jen"
)

const formatSourceCode = true

func WriteGo(statement *jen.Statement) (string, error) {
	if !formatSourceCode {
		return statement.GoString(), nil
	}

	src, err := format.Source([]byte(strings.TrimSpace(statement.GoString())))
	if err != nil {
		return "", fmt.Errorf("Couldn't format source code 2: %w", err)
	}
	return string(src), nil
}

func WriteGoFile(f *jen.File) (string, error) {
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
