package procast

type (
	ErrorHandler func(err error)
)

var EHEmpty ErrorHandler = func(error) {}

func GetRewriteErrHandler(errPtr *error) ErrorHandler {
	if errPtr == nil {
		return EHEmpty
	}
	return func(e error) {
		*errPtr = e
	}
}

func GetWrapVoidErrHandler(fn func()) ErrorHandler {
	if fn == nil {
		return EHEmpty
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
