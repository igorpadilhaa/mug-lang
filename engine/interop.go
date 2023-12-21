package engine

var variables map[string]MugValue = map[string]MugValue{}

var functions map[string]MugFunc = map[string]MugFunc{}

func SetVar(name string, value interface{}) error {
	mValue, err := NewValue(value)
	if err != nil {
		return err
	}

	variables[name] = mValue
	return nil
}

func DefineFunc(name string, fn MugFunc) {
	functions[name] = fn
}

type CallContext struct {
	Args     []MugValue
	err      error
	retValue MugValue
}

func (ctx *CallContext) Error(err error) {
	ctx.err = err
}

func (ctx *CallContext) Return(value MugValue) {
	ctx.retValue = value
}

type MugFunc func(ctx *CallContext)
