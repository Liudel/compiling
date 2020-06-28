package ast

import "fmt"

// ASTNodeType 节点类型
type ASTNodeType string

const (
	// UnKnow 不知道
	UnKnow ASTNodeType = "unKnow"

	// Programm 程序入口
	Programm ASTNodeType = "Programm"

	// IntDeclaration 整型变量声明
	IntDeclaration ASTNodeType = "IntDeclaration"
	// ExpressionStmt 表达式语句
	ExpressionStmt ASTNodeType = "ExpressionStmt"
	// AssignmentStmt 赋值语句
	AssignmentStmt ASTNodeType = "AssignmentStmt"

	// Primary 基础表达式
	Primary ASTNodeType = "Primary"
	// Multiplicative 乘法表达式
	Multiplicative ASTNodeType = "Multiplicative"
	// Additive 加法表达式
	Additive ASTNodeType = "Additive"

	// Identifier 标识符
	Identifier ASTNodeType = "Identifier"
	// IntLiteral 整型字面量
	IntLiteral ASTNodeType = "IntLiteral"
)

// ASTNode ast节点
type ASTNode interface {
	// GetParent 父节点
	GetParent() ASTNode

	// GetChildren 子节点
	GetChildren() []ASTNode

	// GetType ast类型
	GetType() ASTNodeType

	// GetText 文本值
	GetText() string
}

// SimpleASTNode 简单实现
type SimpleASTNode struct {
	Parent           *SimpleASTNode
	Childeren        []ASTNode
	readonlyChildren []ASTNode
	NodeType         ASTNodeType
	Text             string
}

// GetParent 获取父节点
func (sa *SimpleASTNode) GetParent() ASTNode {
	return sa.Parent
}

// GetChildren 获取子节点
func (sa *SimpleASTNode) GetChildren() []ASTNode {
	return sa.Childeren
}

// GetType ast类型
func (sa *SimpleASTNode) GetType() ASTNodeType {
	return sa.NodeType
}

// GetText 文本值
func (sa *SimpleASTNode) GetText() string {
	return sa.Text
}

// AddChild 添加子节点
func (sa *SimpleASTNode) AddChild(child *SimpleASTNode) {
	sa.Childeren = append(sa.Childeren, child)
	child.Parent = sa
}

// DumpAST 打印语法树
func DumpAST(node ASTNode, indent string) {
	fmt.Println(indent + string(node.GetType()) + " " + node.GetText())
	children := node.GetChildren()
	for i := range children {
		DumpAST(children[i], indent+"\t")
	}
}
