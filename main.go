package main

import (
	"github.com/igorpadilhaa/mug/engine"
	"github.com/igorpadilhaa/mug/parser"
	"github.com/igorpadilhaa/mug/lexer"

	"fmt"
)

func main() {
	err := Run(`
		print(hello, world)
		print()
		`)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

func Run(code string) error {
	tokenizer := lexer.Tokenizer{}

	tokens, err := tokenizer.Tokens(code)
	if err != nil {
		return fmt.Errorf("tokenization error: %w", err)
	}
	fmt.Println(tokens)

    node, err := parser.ParseProgram(tokens)
	if err != nil {
		return err
	}

	val, err := engine.Eval(node)
	if err != nil {
		return err
	}

	fmt.Println(val)
	return nil
}