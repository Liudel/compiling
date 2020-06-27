package ast

// ASTNodeType 节点类型
type ASTNodeType int

const (
	// UnKnow 不知道
	UnKnow ASTNodeType = iota

	// Programm 程序入口
	Programm

	// IntDeclaration 整型变量声明
	IntDeclaration
	// ExpressionStmt 表达式语句
	ExpressionStmt
	// AssignmentStmt 赋值语句
	AssignmentStmt

	// Primary 基础表达式
	Primary
	// Multiplicative 乘法表达式
	Multiplicative
	// Additive 加法表达式
	Additive

	// Identifier 标识符
	Identifier
	// IntLiteral 整型字面量
	IntLiteral
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
