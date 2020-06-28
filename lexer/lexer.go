package lexer

import (
	"bytes"
	"cp/token"
)

// DfaState 状态机的状态
type DfaState int

const (
	UnKnow DfaState = iota
	Initial

	If
	IdIf1
	IdIf2
	Else
	IdElse1
	IdElse2
	IdElse3
	IdElse4
	Int
	IdInt1
	IdInt2
	IdInt3
	Id
	GT
	GE
	Assignment

	Plus
	Minus
	Star
	Slash

	SemiColon
	LeftParen
	RightParen

	IntLiteral
)

// SimpleLexer 简单的词法分析器
type SimpleLexer struct {
	tokenText bytes.Buffer
	tokens    []token.Token
	token     token.SimpleToken
}

// isAlpha 是否为字母
func (sl SimpleLexer) isAlpha(ch byte) bool {
	return (int(ch) >= int('a') && int(ch) <= int('z')) || (int(ch) >= int('A') && int(ch) <= int('Z'))
}

// isDigit 是否为数字
func (sl SimpleLexer) isDigit(ch byte) bool {
	return int(ch) >= int('0') && int(ch) <= int('9')
}

// isBlank 是否为空白符
func (sl SimpleLexer) isBlank(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// InitToken 初始化token
// 有限状态机进入初始状态。
// 这个初始状态其实并不做停留，它马上进入其他状态。
// 开始解析的时候，进入初始状态；某个Token解析完毕，也进入初始状态，在这里把Token记下来，然后建立一个新的Token。
func (sl *SimpleLexer) InitToken(ch byte) DfaState {
	if sl.tokenText.Len() > 0 {
		sl.token.Text = sl.tokenText.String()
		sl.tokens = append(sl.tokens, sl.token)

		sl.tokenText.Reset()
		sl.token = token.SimpleToken{}
	}

	state := Initial
	if sl.isAlpha(ch) {
		if ch == 'i' {
			state = IdInt1
		} else {
			state = Id //进入Id状态
		}
		sl.token.Ttype = token.Identifier
	} else if sl.isDigit(ch) {
		state = IntLiteral
		sl.token.Ttype = token.IntLiteral
	} else if ch == '>' {
		state = GT
		sl.token.Ttype = token.GT
	} else if ch == '+' {
		state = Plus
		sl.token.Ttype = token.Plus
	} else if ch == '-' {
		state = Minus
		sl.token.Ttype = token.Minus
	} else if ch == '*' {
		state = Star
		sl.token.Ttype = token.Star
	} else if ch == '/' {
		state = Slash
		sl.token.Ttype = token.Slash
	} else if ch == ';' {
		state = SemiColon
		sl.token.Ttype = token.SemiColon
	} else if ch == '(' {
		state = LeftParen
		sl.token.Ttype = token.LeftParen
	} else if ch == ')' {
		state = RightParen
		sl.token.Ttype = token.RightParen
	} else if ch == '=' {
		state = Assignment
		sl.token.Ttype = token.Assignment
	} else {
		return Initial
	}

	sl.tokenText.WriteByte(ch)
	return state
}

// Tokenize 这是一个有限状态自动机，在不同的状态中迁移
func (sl *SimpleLexer) Tokenize(code string) token.TokenReader {
	sl.tokens = make([]token.Token, 0)
	sl.tokenText = bytes.Buffer{}
	sl.token = token.SimpleToken{}

	state := Initial
	var ch byte
	for i := 0; i < len(code); i++ {
		ch = code[i]
		switch state {
		case Initial:
			state = sl.InitToken(ch)
		case Id:
			if sl.isAlpha(ch) || sl.isDigit(ch) {
				sl.tokenText.WriteByte(ch) //保持标识符状态
			} else {
				state = sl.InitToken(ch) //退出标识符状态，并保存Token
			}
		case GT:
			if ch == '=' {
				sl.token.Ttype = token.GE //转换成GE
				state = GE
				sl.tokenText.WriteByte(ch)
			} else {
				state = sl.InitToken(ch) //退出GT状态，并保存Token
			}
		case GE:
			fallthrough
		case Assignment:
			fallthrough
		case Plus:
			fallthrough
		case Minus:
			fallthrough
		case Star:
			fallthrough
		case Slash:
			fallthrough
		case SemiColon:
			fallthrough
		case LeftParen:
			fallthrough
		case RightParen:
			state = sl.InitToken(ch) //退出当前状态，并保存Token
		case IntLiteral:
			if sl.isDigit(ch) {
				sl.tokenText.WriteByte(ch) //继续保持在数字字面量状态
			} else {
				state = sl.InitToken(ch) //退出当前状态，并保存Token
			}
		case IdInt1:
			if ch == 'n' {
				state = IdInt2
				sl.tokenText.WriteByte(ch)
			} else if sl.isDigit(ch) || sl.isAlpha(ch) {
				state = Id //切换回Id状态
				sl.tokenText.WriteByte(ch)
			} else {
				state = sl.InitToken(ch)
			}
		case IdInt2:
			if ch == 't' {
				state = IdInt3
				sl.tokenText.WriteByte(ch)
			} else if sl.isDigit(ch) || sl.isAlpha(ch) {
				state = Id //切换回id状态
				sl.tokenText.WriteByte(ch)
			} else {
				state = sl.InitToken(ch)
			}
		case IdInt3:
			if sl.isBlank(ch) {
				sl.token.Ttype = token.Int
				state = sl.InitToken(ch)
			} else {
				state = Id //切换回Id状态
				sl.tokenText.WriteByte(ch)
			}
		default:
		}

	}
	if sl.tokenText.Len() > 0 {
		sl.InitToken(ch)
	}

	return &token.SimpleTokenReader{Tokens: sl.tokens}
}
