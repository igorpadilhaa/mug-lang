package engine

import (
	"github.com/igorpadilhaa/mug/lexer"
	"github.com/igorpadilhaa/mug/parser"

	"fmt"
)

type MugValue struct {
	Type MugType
	data interface{}
}

func (value MugValue) AsString() (string, error) {
	if value.Type == MUG_STRING {
		return value.data.(string), nil
	}

	return "", CastError(MUG_STRING, value.Type)
}

func CastError(expected MugType, got MugType) error {
	return fmt.Errorf("can not cast %v to %v", expected, got)
}

var nothing = MugValue{MUG_NOTHING, nil}

func newValue(data interface{}) MugValue {
	switch data.(type) {
	case string:
		return MugValue{MUG_STRING, data}
	default:
		panic(fmt.Errorf("unsupported conversion to MugValue %T", data))
	}
}

var variables map[string]MugValue = map[string]MugValue{
	"hello": newValue("Hello"),
}

var functions map[string]MugFunc = map[string]MugFunc{
	"print": print,
}

func Eval(node parser.ParsedNode) (MugValue, error) {
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

func evalProgram(program parser.ParsedProgram) (MugValue, error) {
	for _, node := range program {
		_, err := Eval(node)
		if err != nil {
			return nothing, err
		}
	}
	return nothing, nil
}

func evalLiteral(literal parser.ParsedLiteral) (MugValue, error) {
	data := literal.Data
	if data.Type != lexer.TOKEN_STRING {
		return nothing, fmt.Errorf("unable to eval literal %s, unknown type %s", data.Content, data.Type)
	}

	unquoted := data.Content[1 : len(data.Content)-1]
	return newValue(unquoted), nil
}

func evalFunction(fc parser.ParsedFunctionCall) (MugValue, error) {
	fn, exists := functions[fc.Name]
	if !exists {
		return nothing, fmt.Errorf("unknown function %q", fc.Name)
	}

	args, err := evalArgumentList(fc.Args)
	if err != nil {
		return nothing, err
	}

	callCtx := CallContext{
		Args: args,
		retValue: nothing,
	}

	fn(callCtx)
	return callCtx.retValue, callCtx.err
}

func evalVariable(variable parser.ParsedVariable) (MugValue, error) {
	value, exist := variables[variable.Name]
	if !exist {
		return nothing, fmt.Errorf("undeclared variable %q", variable.Name)
	}
	return value, nil
}

func evalArgumentList(list parser.ParsedArgumentList) ([]MugValue, error) {
	var values []MugValue

	for _, item := range list.Values {
		value, err := Eval(item)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func evalVariableAssignment(assignment parser.ParsedVariableAssigment) (MugValue, error) {
	value, err := Eval(assignment.Value)
	if err != nil {
		return nothing, err
	}

	variables[assignment.Variable.Name] = value
	return nothing, nil
}
