#!/bin/sh

# Setup:
# Install antlr: pip install antlr4-tools

# To get from private reops (temporary)
# GOPRIVATE="github.com/mattmunz/appkit" go get github.com/mattmunz/appkit

gen() {
  echo "Gen ANTLR code..."

  # antlr4-parse documentation/DesignLanguage.g4 design ./test/data/Test.1.nzsd.txt -tree

  outDir="model/designlanguage/parser"
  mkdir -p "$outDir"
  # Go runtime (once): go get github.com/antlr/antlr4/runtime/Go/antlr/v4@latest
  tempDir=$(mktemp -d)
  antlr4 -Dlanguage=Go -package parser -o "$tempDir" documentation/design/DesignLanguage.g4
  rm -rf "$outDir/" && mv "$tempDir/documentation"/* "$outDir/"
  rm -rf "$tempDir"

  echo "Gen nzsd models..."

  designlanguage gen -d=false -t=go

  # designlanguage --help
}

test() {
  DATA_PATH="$(realpath ./test/data)"

  TEST_DATA_PATH=$DATA_PATH go test  ./test/unit/... # -test.v -run ^TestLoad1$
}

testGen() {
  designlanguage --help
}

case $1 in
  b  | build) go install ./cmd/designlanguage/designlanguage.go;;
  f  | format) go fmt ./*;;  
  g  | gen) gen;;
  h  | help)  echo "build|gen (antlr)|help|test|test_functional";;
  t  | test)  test;;
  tf | test_functional)  go test ./test/functional/...;;
  tg | test_gen) testGen;;
  *) echo "Unknown argument";;
esac
