package typer

import (
	"fmt"
)

func PanicWhenError(err error, fmtOrMsg string, args ...any) {
	if err != nil {
		if fmtOrMsg == "" {
			panic(err)
		}
		panic(fmt.Errorf(fmtOrMsg+", %w", append(args, err)...))
	}
}
