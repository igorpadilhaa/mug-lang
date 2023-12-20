package engine

import (
	"fmt"
)

func print(ctx CallContext) {
	for _, arg := range ctx.Args {
		str, err := arg.AsString()
		if err != nil {
			ctx.Error(err)
		}

		fmt.Printf("%s ", str)
	}
	fmt.Println()
}
