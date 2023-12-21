package engine

import (
	"github.com/igorpadilhaa/mug/engine"

	"fmt"
)

func init() {
	engine.DefineFunc("print", print)
	engine.DefineFunc("string", toString)
}

func print(ctx *engine.CallContext) {
	for _, arg := range ctx.Args {
		str, err := arg.AsString()
		if err != nil {
			ctx.Error(err)
		}

		fmt.Printf("%s ", str)
	}
	fmt.Println()
}

func toString(ctx *engine.CallContext) {
	if len(ctx.Args) != 1 {
		ctx.Error(fmt.Errorf("argument count mismatch; expected 1, got %d", len(ctx.Args)))
		return
	}

	value := ctx.Args[0]
	switch value.Type {
	case engine.MUG_INT:
		data, _ := value.AsInt()
		stringfied := fmt.Sprintf("%d", data)

		mString, err := engine.NewValue(stringfied)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Return(mString)

	case engine.MUG_STRING:
		ctx.Return(value)

	case engine.MUG_NOTHING:
		mString, err := engine.NewValue("nothing")
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Return(mString)
	}
}
