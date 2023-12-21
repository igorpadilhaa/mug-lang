package lexer

import (
	"fmt"
)

type TokenType int

type Token struct {
	Type    TokenType
	Content string
}

const (
	TOKEN_UNKNOWN TokenType = iota
	TOKEN_SPACE
	TOKEN_NAME
	TOKEN_PAROP
	TOKEN_PARCL
	TOKEN_COMMA
	TOKEN_OPASSING
	TOKEN_STRING
	TOKEN_INTEGER
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case TOKEN_UNKNOWN:
		return "TK_UNKNOWN"
	case TOKEN_SPACE:
		return "TK_SPACE"
	case TOKEN_NAME:
		return "TK_NAME"
	case TOKEN_PAROP:
		return "TK_OPPAR"
	case TOKEN_PARCL:
		return "TK_CLPAR"
	case TOKEN_COMMA:
		return "TK_COMMA"
	case TOKEN_OPASSING:
		return "TK_OP_EQU"
	case TOKEN_STRING:
		return "TK_STRING"
	case TOKEN_INTEGER:
		return "TK_INTEGER"
	default:
		panic(fmt.Errorf("unknown token type %d", tokenType))
	}
}
