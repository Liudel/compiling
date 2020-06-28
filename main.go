package main

import (
	"cp/ast"
	"cp/parser"
)

func main() {
	a := &parser.SimpleParser{}

	ast.DumpAST(a.Parser("int age = 45*2+2;"), "")
}
