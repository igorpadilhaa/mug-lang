package parser

import (
	"github.com/igorpadilhaa/mug/lexer"

	"errors"
	"fmt"
)

var errEOF error = errors.New("unexpected end of file")

type ParsedNode interface {
}

type ParseFunction func(tokenizer *lexer.Tokenizer) (ParsedNode, error)

type ParsedProgram []ParsedNode

func ParseProgram(tokens []lexer.Token) (ParsedProgram, error) {
	tokenizer := lexer.Tokenizer{Content: tokens}

	parsers := []ParseFunction{
		ParseFunctionCall,
		ParseVariableAssignment,
		ParseValue,
	}

	var program ParsedProgram
	for tokenizer.HasNext() {
		node, err := ParseCompoundRule(&tokenizer, parsers)
		if err != nil {
			return nil, err
		}
		program = append(program, node)
	}

	return program, nil
}

func ParseCompoundRule(tokenizer *lexer.Tokenizer, parsers []ParseFunction) (ParsedNode, error) {
	for _, parser := range parsers {
		branch := *tokenizer
		if node, err := parser(&branch); err == nil {
			*tokenizer = branch
			return node, nil
		}
	}
	return nil, fmt.Errorf("failed to parse value; expected %v", parsers)
}

type ParsedFunctionCall struct {
	Name string
	Args ParsedArgumentList
}

func ParseFunctionCall(tokens *lexer.Tokenizer) (ParsedNode, error) {
	var funcCall ParsedFunctionCall

	tk, err := tokens.Next()
	if err != nil {
		return nil, err
	} else if tk.Type != lexer.TOKEN_NAME {
		return nil, ParserError(RULE_FUNCALL, lexer.TOKEN_NAME, tk.Type)
	}
	funcCall.Name = tk.Content

	tk, err = tokens.Next()
	if err != nil {
		return nil, err
	} else if tk.Type != lexer.TOKEN_PAROP {
		return nil, ParserError(RULE_FUNCALL, lexer.TOKEN_NAME, tk.Type)
	}

	argList, err := ParseArgumentList(tokens)
	funcCall.Args = argList

	if err != nil {
		return nil, err
	}

	tk, err = tokens.Next()
	if err != nil {
		return nil, err

	} else if tk.Type != lexer.TOKEN_PARCL {
		return nil, ParserError(RULE_FUNCALL, lexer.TOKEN_PARCL, tk.Type)
	}

	return funcCall, nil
}

type ParsedArgumentList struct {
	Values []ParsedNode
}

func ParseArgumentList(tokens *lexer.Tokenizer) (ParsedArgumentList, error) {
	mayBeEmpty := true
	firstItem := true

	hasComma := false

	var list ParsedArgumentList

	for tokens.HasNext() {
		tk, _ := tokens.Peek()

		branch := *tokens
		node, err := ParseValue(&branch)

		if mayBeEmpty && err != nil {
			return list, nil
		} else {
			mayBeEmpty = false
		}

		if tk.Type == lexer.TOKEN_COMMA {
			if hasComma {
				tokens.Back()
				return list, nil
			}

			hasComma = true
			tokens.Next()
			continue
		}

		if err == nil {
			*tokens = branch
			if !firstItem && !hasComma {
				return list, nil
			}

			firstItem = false
			hasComma = false
			list.Values = append(list.Values, node)

		} else {
			if hasComma {
				tokens.Back()
			}
			return list, nil
		}
	}

	return list, errEOF
}

type ParsedLiteral struct {
	Data lexer.Token
}

func ParseValue(tokenizer *lexer.Tokenizer) (ParsedNode, error) {
	tk, err := tokenizer.Peek()
	if err != nil {
		return nil, errEOF
	}

	switch tk.Type {
	case lexer.TOKEN_STRING, lexer.TOKEN_INTEGER:
		tokenizer.Next()
		return ParsedLiteral{tk}, nil
	}

	parsers := []ParseFunction{
		ParseFunctionCall,
		ParseVariable,
	}

	return ParseCompoundRule(tokenizer, parsers)
}

type ParsedVariable struct {
	Name string
}

func ParseVariable(tokenizer *lexer.Tokenizer) (ParsedNode, error) {
	if !tokenizer.HasNext() {
		return ParsedVariable{}, errEOF
	}

	tk, _ := tokenizer.Peek()
	if tk.Type != lexer.TOKEN_NAME {
		return ParsedVariable{}, ParserError(RULE_VARIABLE, lexer.TOKEN_NAME, tk.Type)
	}

	tokenizer.Next()
	return ParsedVariable{tk.Content}, nil
}

type ParsedVariableAssigment struct {
	Variable ParsedVariable
	Value    ParsedNode
}

func ParseVariableAssignment(tokenizer *lexer.Tokenizer) (ParsedNode, error) {
	var assignment ParsedVariableAssigment

	variable, err := ParseVariable(tokenizer)
	if err != nil {
		return nil, err
	}
	assignment.Variable = variable.(ParsedVariable)

	if tk, err := tokenizer.Next(); err != nil {
		return nil, err
	} else if tk.Type != lexer.TOKEN_OPASSING {
		return nil, ParserError(RULE_VARIABLE_ASSIGNMENT, lexer.TOKEN_OPASSING, tk.Type)
	}

	value, err := ParseValue(tokenizer)
	if err != nil {
		return nil, err
	}
	assignment.Value = value

	return assignment, nil
}

func ParserError(rule RuleType, expected lexer.TokenType, got lexer.TokenType) error {
	return fmt.Errorf("failed to parse %s; expected %s, got %s", rule, expected, got)
}
