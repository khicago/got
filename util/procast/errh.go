package procast

import "github.com/khicago/got/util/delegate"

type (
	ErrorHandler delegate.Action1[error]
)

var EHDoNothing ErrorHandler = func(error) {}

func GetRewriteErrHandler(errPtr *error) ErrorHandler {
	if errPtr == nil {
		return EHDoNothing
	}
	return func(e error) {
		if e != nil {
			*errPtr = e
		}
	}
}

func GetWrapVoidErrHandler(fn delegate.Action) ErrorHandler {
	if fn == nil {
		return EHDoNothing
	}
	return func(e error) {
		fn()
	}
}

func GetCombinedErrHandler(handlers ...ErrorHandler) ErrorHandler {
	return func(e error) {
		for _, h := range handlers {
			if h != nil {
				h(e)
			}
		}
	}
}

func (eg ErrorHandler) SafeGo(fn func()) {
	SafeGo(fn, eg)
}
