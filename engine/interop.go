package engine

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

type MugFunc func(ctx CallContext)
