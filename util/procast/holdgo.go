package procast

import (
	"sync"

	"github.com/khicago/got/util/delegate"
)

// AtLeastOnce for returning a pair of functions.
// The `wait` method can block a goroutine, and the `exit` method causes
// this blocking to exit.
// `skipHolder` can be called before or after the `wait`,
// It also can be called multiple times, but only the first call makes sense.
func AtLeastOnce() (wait delegate.Action, skipHolder delegate.Action) {
	wg, once := sync.WaitGroup{}, sync.Once{}
	wg.Add(1)
	return wg.Wait, func() {
		once.Do(wg.Done)
	}
}

// HoldGo - Hold the proc until closer are called or panic
// closer can be called multi-times
// the quiting of HoldGo dose not means that the fn are finished
func HoldGo(fn func(skip ErrorHandler)) (err error) {
	wait, exitHolder := AtLeastOnce()
	handleErrAndExit := GetCombinedErrHandler(GetRewriteErrHandler(&err), GetWrapVoidErrHandler(exitHolder))
	handleErrAndExit.SafeGo(func() {
		defer exitHolder()
		fn(handleErrAndExit)
	})
	wait() // hold go
	return
}
