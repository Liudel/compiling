package parser

import (
	"cp/ast"
	"cp/lexer"
	"cp/token"
)

// SimpleParser 简单的语法解析器
type SimpleParser struct{}

// Parser 解析脚本
func (sp *SimpleParser) Parser(script string) ast.ASTNode {
	lexer := &lexer.SimpleLexer{}
	tokens := lexer.Tokenize(script)
	rootNode := sp.Prog(tokens)
	return rootNode
}

// Prog 解析入口
func (sp *SimpleParser) Prog(tokens token.TokenReader) ast.ASTNode {
	node := &ast.SimpleASTNode{NodeType: ast.Programm, Text: "pwc"}

	for tokens.Peek() != nil {
		child := sp.IntDeclare(tokens)
		if child == nil {
			child = sp.ExpressionStatement(tokens)
		}

		if child == nil {
			child = sp.AssignmentStatement(tokens)
		}

		if child != nil {
			node.AddChild(child)
		} else {
			panic("unknown statement")
		}
	}

	return node
}

// IntDeclare 整型变量声明
func (sp *SimpleParser) IntDeclare(tokens token.TokenReader) *ast.SimpleASTNode {
	var node *ast.SimpleASTNode
	oneToken := tokens.Peek()
	if oneToken != nil && oneToken.GetType() == token.Int {
		tokens.Read()
		if tokens.Peek().GetType() == token.Identifier {
			oneToken = tokens.Read()
			node = &ast.SimpleASTNode{NodeType: ast.IntDeclaration, Text: oneToken.GetText()}
			oneToken = tokens.Peek()
			if oneToken != nil && oneToken.GetType() == token.Assignment {
				tokens.Read()
				child := sp.Additive(tokens)
				if child == nil {
					panic("invalide variable initialization, expecting an expression")
				}
				node.AddChild(child)
			}
		} else {
			panic("variable name expected")
		}

		if node != nil {
			oneToken = tokens.Peek() //预读，看看后面是不是分号
			if oneToken != nil && oneToken.GetType() == token.SemiColon {
				tokens.Read() // 消耗掉这个分号
			} else {
				panic("invalid statement, expecting semicolon")
			}
		}
	}
	return node
}

// ExpressionStatement 表达式语句，即表达式后面跟个分号
func (sp *SimpleParser) ExpressionStatement(tokens token.TokenReader) *ast.SimpleASTNode {
	pos := tokens.GetPosition()
	node := sp.Additive(tokens)
	if node != nil {
		oneToken := tokens.Peek()
		if oneToken != nil && oneToken.GetType() == token.SemiColon {
			tokens.Read()
		} else {
			node = nil
			tokens.SetPosition(pos) // 回溯
		}
	}
	return node // 直接返回子节点
}

// AssignmentStatement 赋值语句
func (sp *SimpleParser) AssignmentStatement(tokens token.TokenReader) *ast.SimpleASTNode {
	var node *ast.SimpleASTNode
	oneToken := tokens.Peek()
	if oneToken != nil && oneToken.GetType() == token.Identifier {
		oneToken = tokens.Read()
		node = &ast.SimpleASTNode{NodeType: ast.AssignmentStmt, Text: oneToken.GetText()}
		oneToken = tokens.Peek()
		if oneToken != nil && oneToken.GetType() == token.Assignment {
			tokens.Read()
			child := sp.Additive(tokens)
			if child == nil {
				panic("invalide assignment statement, expecting an expression")
			}

			node.AddChild(child)
			oneToken = tokens.Peek() //预读，看看后面是不是分号
			if oneToken != nil && oneToken.GetType() == token.SemiColon {
				tokens.Read() // 消耗掉这个分号
			} else {
				panic("invalid statement, expecting semicolon")
			}
		} else {
			tokens.Unread() //回溯，吐出之前消化掉的标识符
			node = nil
		}
	}

	return node
}

// Additive 加法表达式
func (sp *SimpleParser) Additive(tokens token.TokenReader) *ast.SimpleASTNode {
	child1 := sp.Multiplicative(tokens)
	node := child1

	if child1 != nil {
		for {
			oneToken := tokens.Peek()
			if oneToken != nil && (oneToken.GetType() == token.Plus || oneToken.GetType() == token.Minus) {
				oneToken = tokens.Read()
				child2 := sp.Multiplicative(tokens)
				if child2 == nil {
					panic("invalid additive expression, expecting the right part.")
				}

				node = &ast.SimpleASTNode{NodeType: ast.Additive, Text: oneToken.GetText()}
				node.AddChild(child1)
				node.AddChild(child2)
				child1 = node

			} else {
				break
			}
		}
	}
	return node
}

// AdditiveOne 加法表达式结合错误
func (sp *SimpleParser) AdditiveOne(tokens token.TokenReader) *ast.SimpleASTNode {
	child1 := sp.Multiplicative(tokens)
	node := child1
	oneToken := tokens.Peek()
	if child1 != nil && oneToken != nil {
		if oneToken.GetType() == token.Plus || oneToken.GetType() == token.Minus {
			oneToken = tokens.Read()
			child2 := sp.Additive(tokens)
			if child2 == nil {
				panic("invalid additive expression, expecting the right part.")
			}

			node = &ast.SimpleASTNode{NodeType: ast.Additive, Text: oneToken.GetText()}
			node.AddChild(child1)
			node.AddChild(child2)
			child1 = node

		}
	}

	return node
}

// Multiplicative 乘法表达式
func (sp *SimpleParser) Multiplicative(tokens token.TokenReader) *ast.SimpleASTNode {
	child1 := sp.Primary(tokens)
	node := child1

	for {
		oneToken := tokens.Peek()
		if oneToken != nil && (oneToken.GetType() == token.Star || oneToken.GetType() == token.Slash) {
			oneToken = tokens.Read()
			child2 := sp.Primary(tokens)
			if child2 == nil {
				panic("invalid multiplicative expression, expecting the right part.")
			}

			node = &ast.SimpleASTNode{NodeType: ast.Multiplicative, Text: oneToken.GetText()}
			node.AddChild(child1)
			node.AddChild(child2)
			child1 = node
		} else {
			break
		}
	}
	return node
}

// Primary 基础表达式
func (sp *SimpleParser) Primary(tokens token.TokenReader) *ast.SimpleASTNode {
	var node *ast.SimpleASTNode
	oneToken := tokens.Peek()
	if tokens != nil {
		if oneToken.GetType() == token.IntLiteral {
			oneToken = tokens.Read()
			node = &ast.SimpleASTNode{NodeType: ast.IntLiteral, Text: oneToken.GetText()}
		} else if oneToken.GetType() == token.Identifier {
			oneToken = tokens.Read()
			node = &ast.SimpleASTNode{NodeType: ast.Identifier, Text: oneToken.GetText()}
		} else if oneToken.GetType() == token.LeftParen {
			tokens.Read()
			node = sp.Additive(tokens)
			if node != nil {
				oneToken = tokens.Peek()
				if oneToken != nil && oneToken.GetType() == token.RightParen {
					tokens.Read()
				} else {
					panic("expecting right parenthesis")
				}
			} else {
				panic("expecting an additive expression inside parenthesis")
			}
		}
	}
	return node
}
