package lexer

import (
	"errors"
	"fmt"
	"unicode"
)

type Tokenizer struct {
	Position int
	Content  []Token
}

func (ti *Tokenizer) HasNext() bool {
	return ti.Position < len(ti.Content)
}

func (ti *Tokenizer) Peek() (Token, error) {
	if ti.HasNext() {
		index := ti.Position
		return ti.Content[index], nil
	}
	return Token{}, errors.New("reached end of list; failed to retrieve next token")
}

func (ti *Tokenizer) Next() (Token, error) {
	token, err := ti.Peek()
	if err == nil {
		ti.Position += 1
	}
	return token, err
}

func (ti *Tokenizer) Back() (Token, error) {
	if ti.Position == 0 {
		return Token{}, errors.New("no more elements; failed to move to back element")
	}
	index := ti.Position
	ti.Position -= 1
	return ti.Content[index], nil
}

func (tokenizer *Tokenizer) Tokens(code string) ([]Token, error) {
	var tokens []Token
	var blob string

	tokenType := TOKEN_UNKNOWN

	appendToken := false
	runes := []rune(code)

	for tokenizer.Position < len(runes) {
		letter := runes[tokenizer.Position]

		switch tokenType {
		case TOKEN_UNKNOWN:
			switch {
			case unicode.IsSpace(letter):
				tokenType = TOKEN_SPACE

			case letter == '"':
				tokenType = TOKEN_STRING

			case letter == '=':
				tokenType = TOKEN_OPASSING

			case letter == '(':
				tokenType = TOKEN_PAROP

			case letter == ')':
				tokenType = TOKEN_PARCL

			case letter == ',':
				tokenType = TOKEN_COMMA

			case unicode.IsLetter(letter):
				tokenType = TOKEN_NAME

			default:
				return nil, tokenizerError(tokenizer.Position, letter)
			}

		case TOKEN_SPACE:
			if !unicode.IsSpace(letter) {
				tokenType = TOKEN_UNKNOWN
				break
			}

			tokenizer.Position += 1

		case TOKEN_PAROP, TOKEN_PARCL, TOKEN_COMMA, TOKEN_OPASSING:
			blob += string(letter)
			appendToken = true
			tokenizer.Position += 1

		case TOKEN_NAME:
			if !unicode.IsLetter(letter) {
				appendToken = true
				break
			}

			blob += string(letter)
			tokenizer.Position += 1

		case TOKEN_STRING:
			if len(blob) == 0 && letter != '"' {
				return nil, tokenizerError(tokenizer.Position, letter, TOKEN_STRING)
			}
			
			blob += string(letter)
			tokenizer.Position += 1

			if len(blob) != 1 && letter == '"' {
				appendToken = true
			}
		}

		if appendToken {
			tokens = append(tokens, Token{tokenType, blob})
			blob = ""
			tokenType = TOKEN_UNKNOWN
			appendToken = false
		}
	}

	if len(blob) != 0 {
		tokens = append(tokens, Token{tokenType, blob})
	}

	return tokens, nil
}

func tokenizerError(position int, got rune, expected ...TokenType) error {
	if len(expected) != 0 {
		return fmt.Errorf("unknown token starting with %q at %d; expected %q", got, position, expected)
	}

	return fmt.Errorf("unknown token starting with %q at %d", got, position)
}
