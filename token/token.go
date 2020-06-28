package token

import (
	"fmt"
)

// TokenType 词法类型
type TokenType string

const (
	// UnKnow 不知道
	UnKnow TokenType = "UnKnow"
	// Plus 加
	Plus TokenType = "Plus"
	// Minus 减
	Minus TokenType = "Minus"
	// Star 乘
	Star TokenType = "Star"
	// Slash 除
	Slash TokenType = "Slash"

	// GE 大于等于
	GE TokenType = "GE"
	// GT 大于
	GT TokenType = "GT"
	// EQ 等于
	EQ TokenType = "GQ"
	// LE 小于等于
	LE TokenType = "LE"
	// LT 小于
	LT TokenType = "LT"

	// SemiColon 分号
	SemiColon TokenType = "SemiColon"
	// LeftParen 左括号
	LeftParen TokenType = "LeftParen"
	// RightParen 右括号
	RightParen TokenType = "RightParen"

	// Assignment 赋值
	Assignment TokenType = "Assignment"

	// If if
	If TokenType = "If"
	// Else else
	Else TokenType = "Else"

	// Int int
	Int TokenType = "Int"

	// Identifier 标识符
	Identifier TokenType = "Identifier"

	// IntLiteral 整型字面量
	IntLiteral TokenType = "IntLiteral"
	// StringLiteral 字符串字面量
	StringLiteral TokenType = "StringLiteral"
)

// Token 词法分析接口
type Token interface {
	GetType() TokenType

	GetText() string
}

// TokenReader token流
type TokenReader interface {
	// 返回Token流中下一个Token，并从流中取出。 如果流已经为空，返回null;
	Read() Token
	// 返回Token流中下一个Token，但不从流中取出。 如果流已经为空，返回null;
	Peek() Token

	// Token流回退一步。恢复原来的Token。
	Unread()

	// 获取Token流当前的读取位置。
	GetPosition() int
	// 设置Token流当前的读取位置
	SetPosition(int)
}

// SimpleToken 简单的token实现
type SimpleToken struct {
	Ttype TokenType
	Text  string
}

// GetText 获取文本
func (st SimpleToken) GetText() string {
	return st.Text
}

// GetType 获取类型
func (st SimpleToken) GetType() TokenType {
	return st.Ttype
}

// SimpleTokenReader 词法读取器
type SimpleTokenReader struct {
	Tokens []Token
	pos    int
}

// Read 读取
func (s *SimpleTokenReader) Read() Token {
	if s.pos < len(s.Tokens) {
		result := s.Tokens[s.pos]
		s.pos++
		return result
	}
	return nil
}

// Peek 选取
func (s *SimpleTokenReader) Peek() Token {
	if s.pos < len(s.Tokens) {
		return s.Tokens[s.pos]
	}
	return nil
}

// Unread 回退
func (s *SimpleTokenReader) Unread() {
	if s.pos > 0 {
		s.pos--
	}
}

// GetPosition 获取位置
func (s *SimpleTokenReader) GetPosition() int {
	return s.pos
}

// SetPosition 设置位置
func (s *SimpleTokenReader) SetPosition(position int) {
	if position >= 0 && position < len(s.Tokens) {
		s.pos = position
	}
}

// Dump 打印词法规则
func Dump(reader TokenReader) {
	fmt.Println("text\ttype")
	var token Token

	for token = reader.Read(); token != nil; token = reader.Read() {
		fmt.Println(token.GetText() + "\t\t" + string(token.GetType()))
	}
}
