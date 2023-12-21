package engine

import (
	"github.com/igorpadilhaa/mug/engine"
	
	"fmt"
)

func init() {
	engine.DefineFunc("print", print)
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
