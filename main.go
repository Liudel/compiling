package main

import (
	"cp/ast"
	"cp/parser"
)

func main() {
	a := &parser.SimpleParser{}

	ast.DumpAST(a.Parser("2+3+4;"), "")
}
