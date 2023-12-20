package engine

import (
	"github.com/igorpadilhaa/mug/parser"
	"github.com/igorpadilhaa/mug/lexer"

	"fmt"
)

type Value interface {
}

var variables map[string]Value = map[string]Value{
	"hello": "Hello",
}

func Eval(node parser.ParsedNode) (Value, error) {
	switch t := node.(type) {
	case parser.ParsedProgram:
		return evalProgram(t)

	case parser.ParsedLiteral:
		return evalLiteral(t)

	case parser.ParsedFunctionCall:
		return evalFunction(t)

	case parser.ParsedVariable:
		return evalVariable(t)

	case parser.ParsedVariableAssigment:
		return evalVariableAssignment(t)

	default:
		panic(fmt.Sprintf("unsupported node %#v", node))
	}
}

func evalProgram(program parser.ParsedProgram) (Value, error) {
	for _, node := range program {
		_, err := Eval(node)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func evalLiteral(literal parser.ParsedLiteral) (Value, error) {
	data := literal.Data
	if data.Type != lexer.TOKEN_STRING {
		return nil, fmt.Errorf("unable to eval literal %s, unknown type %s", data.Content, data.Type)
	}

	unquoted := data.Content[1:len(data.Content)-1]
	return unquoted, nil
}

func evalFunction(fc parser.ParsedFunctionCall) (Value, error) {
	if fc.Name != "print" {
		return nil, fmt.Errorf("unknown function %q", fc.Name)
	}

	args, err := evalArgumentList(fc.Args)
	if err != nil {
		return nil, err
	}

	var anys []any
	for _, arg := range args {
		anys = append(anys, arg)
	}

	fmt.Println(anys...)
	return nil, nil
}

func evalVariable(variable parser.ParsedVariable) (Value, error) {
	value, exist := variables[variable.Name]
	if !exist {
		return nil, fmt.Errorf("undeclared variable %q", variable.Name)
	}
	return value, nil
}

func evalArgumentList(list parser.ParsedArgumentList) ([]Value, error) {
	var values []Value
	
	for _, item := range list.Values {
		value, err := Eval(item)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func evalVariableAssignment(assignment parser.ParsedVariableAssigment) (Value, error) {
	value, err := Eval(assignment.Value)
	if err != nil {
		return nil, err
	}

	variables[assignment.Variable.Name] = value 
	return nil, nil
}