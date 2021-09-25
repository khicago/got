package procast

type (
	ErrorHandler func(err error)
)

var EHEmpty ErrorHandler = func(error) {}

func GetRewriteErrHandler(err *error) ErrorHandler {
	if err == nil {
		return EHEmpty
	}
	return func(e error) {
		*err = e
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

func (eg ErrorHandler) Recover() {
	innerRecover(eg)
}

func (eg ErrorHandler) SafeGo(fn func()) {
	SafeGo(fn, eg)
}
