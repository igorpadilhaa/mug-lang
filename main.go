package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/igorpadilhaa/mug/engine"
	"github.com/igorpadilhaa/mug/lexer"
	"github.com/igorpadilhaa/mug/parser"
	_ "github.com/igorpadilhaa/mug/std"

	"fmt"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	repl(reader)
}

func repl(reader *bufio.Reader) {
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ERROR: unexpected end of line")
			return			
		}

		line = strings.TrimSpace(line)
		if line == ".quit" {
			fmt.Println("Goodbye!")
			return
		}

		err = Run(line)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}
}

func Run(code string) error {
	tokenizer := lexer.Tokenizer{}

	tokens, err := tokenizer.Tokens(code)
	if err != nil {
		return fmt.Errorf("tokenization error: %w", err)
	}

    node, err := parser.ParseProgram(tokens)
	if err != nil {
		return err
	}

	_, err = engine.Eval(node)
	if err != nil {
		return err
	}

	return nil
}