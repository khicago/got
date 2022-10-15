package delegate

type (
	Func[TOut any]                  func() TOut
	Func1[TIn any, TOut any]        func(TIn) TOut
	Func2[TIn1, TIn2 any, TOut any] func(TIn1, TIn2) TOut
)

// Partial
// partial func1 to func
func ToFunc1[TIn any, TOut any](f Func1[TIn, TOut]) Func1[TIn, TOut] {
	return f
}

// Partial
// partial func1 to func
func (t Func1[TIn, TOut]) Partial(val TIn) Func[TOut] {
	return func() TOut {
		return t(val)
	}
}

// TryCall
// do nothing when fulc is nil
func (t Func1[TIn, TOut]) TryCall(val TIn, defaultVal TOut) TOut {
	if t == nil {
		return defaultVal
	}
	return t(val)
}

// Partial
// partial func2 to func1
func (t Func2[TIn1, TIn2, TOut]) Partial(val TIn1) Func1[TIn2, TOut] {
	return func(v2 TIn2) TOut {
		return t(val, v2)
	}
}
