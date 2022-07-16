package inlog

import (
	"fmt"
	"github.com/khicago/got/util/typer"
)

type (
	IMinimLog interface {
		Debugf(format string, args ...any)
		Infof(format string, args ...any)
		Warnf(format string, args ...any)
		Errorf(format string, args ...any)
		Panicf(format string, args ...any)

		Debug(args ...any)
		Info(args ...any)
		Warn(args ...any)
		Error(args ...any)
		Panic(args ...any)
	}
)

func print(prefix, format string, args ...any) {
	if !typer.IsZero(format) {
		fmt.Printf(prefix+" "+format, args...)
	}
	fmt.Println(append([]any{prefix}, args...)...)
}
