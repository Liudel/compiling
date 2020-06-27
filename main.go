package main

import (
	"cp/lexer"
)

func main() {
	a := &lexer.SimpleLexer{}
	lexer.Dump(a.Tokenize("2+3*5"))
}
