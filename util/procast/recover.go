package procast

import (
	"fmt"
	"runtime"
	"strings"
)

func figurePanicLocation(skipExtra int) string {
	var frame runtime.Frame
	var fnName string

	pc := make([]uintptr, 16)
	_ = runtime.Callers(3+skipExtra, pc)
	frames := runtime.CallersFrames(pc)

	for i := 0; i < len(pc); i++ {
		if frames == nil {
			break
		}

		f, more := frames.Next()
		fn := runtime.FuncForPC(frame.PC)
		if fn == nil {
			break
		}
		frame = f
		fnName = fn.Name()
		if !strings.HasPrefix(fnName, "runtime.") {
			break
		}

		if !more {
			break
		}
	}

	if fnName != "" {
		return fmt.Sprintf("[func] %v:%v", fnName, frame.Line)
	}

	if frame.File != "" {
		return fmt.Sprintf("[file] %v:%v", frame.File, frame.Line)
	}

	return fmt.Sprintf("[pc] %x", pc)
}

// Recover 适用于需要对 err 进行重写的场景, 如 handler 中
func Recover(handler ErrorHandler) {
	innerRecover(handler)
}

// Recover 适用于需要对 err 进行重写的场景, 如 handler 中
func innerRecover(handler ErrorHandler) {
	e := recover()
	if handler == nil || e == nil {
		return
	}
	panicLoc := figurePanicLocation(1)
	if err, ok := e.(error); ok {
		handler(fmt.Errorf("%w [ panic !!! %s ]", err, panicLoc))
	} else {
		handler(fmt.Errorf("%v [ panic !!! %s ]", e, panicLoc))
	}
}

func SafeGo(fn func(), errHandler ErrorHandler) {
	go func() {
		defer Recover(errHandler)
		fn()
	}()
	return
}
