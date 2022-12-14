package procast

import (
	"fmt"
	"runtime/debug"
)

func Recover(handler ErrorHandler) {
	e := recover()
	if handler == nil || e == nil {
		return
	}
	panicLoc := figurePanicLocation()
	if err, ok := e.(error); ok {
		handler(fmt.Errorf("%w [ panic !!! %s ]", err, panicLoc))
	} else {
		handler(fmt.Errorf("%v [ panic !!! %s ]", e, panicLoc))
	}
	// innerRecover(handler)
}

func innerRecover(handler ErrorHandler) {
	e := recover()
	if handler == nil || e == nil {
		return
	}
	panicLoc := figurePanicLocation()
	stack := debug.Stack()
	// fmt.Println("stack", string(stack))
	if err, ok := e.(error); ok {
		handler(fmt.Errorf("%w [ panic !!! %s , stack= %s ]", err, panicLoc, string(stack)))
	} else {
		handler(fmt.Errorf("%v [ panic !!! %s , stack= %s ]", e, panicLoc, string(stack)))
	}
}

// SafeGo run function in a protected goroutine
func SafeGo(fn func(), panicErrHandler ErrorHandler) {
	go func() {
		defer innerRecover(panicErrHandler)
		fn()
	}()
}

func figurePanicLocation() string {
	fNode := GetFrameNode(3, "runtime.")
	return fNode.LocString()
}
