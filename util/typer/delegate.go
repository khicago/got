package typer

type (
	DelegateAction                      func()
	DelegateAction1[T any]              func(T)
	DelegateAction2[T1, T2 any]         func(T1, T2)
	DelegateAction3[T1, T2, T3 any]     func(T1, T2, T3)
	DelegateAction4[T1, T2, T3, T4 any] func(T1, T2, T3, T4)

	DelegateFunc[TOut any]                  func() TOut
	DelegateFunc1[TIn any, TOut any]        func(TIn) TOut
	DelegateFunc2[TIn1, TIn2 any, TOut any] func(TIn1, TIn2) TOut

	DelegateHandler[TError error]                  DelegateFunc[TError]
	DelegateHandler1[TIn any, TError error]        DelegateFunc1[TIn, TError]
	DelegateHandler2[TIn1, TIn2 any, TError error] DelegateFunc2[TIn1, TIn2, TError]

	DelegateMap[TIn any, TOut any] DelegateFunc1[TIn, TOut]

	Predicate[TVal any] DelegateMap[TVal, bool]
	Compare[TVal any]   func(TVal, TVal) int
)

// Partial
// partial func1 to func
func (t DelegateAction1[T]) Partial(val T) DelegateAction {
	return func() {
		t(val)
	}
}

// TryCall
// do nothing when action is nil
func (t DelegateAction1[T]) TryCall(val T) {
	if t == nil {
		return
	}
	t(val)
}

// Partial
// partial func1 to func
func (t DelegateAction2[T1, T2]) Partial(v1 T1) DelegateAction1[T2] {
	return func(v2 T2) {
		t(v1, v2)
	}
}

// TryCall
// do nothing when action is nil
func (t DelegateAction2[T1, T2]) TryCall(v1 T1, v2 T2) {
	if t == nil {
		return
	}
	t(v1, v2)
}

// Partial
// partial func1 to func
func (t DelegateAction3[T1, T2, T3]) Partial(v1 T1) DelegateAction2[T2, T3] {
	return func(v2 T2, v3 T3) {
		t(v1, v2, v3)
	}
}

// Partial
// partial func1 to func
func (t DelegateAction4[T1, T2, T3, T4]) Partial(v1 T1) DelegateAction3[T2, T3, T4] {
	return func(v2 T2, v3 T3, v4 T4) {
		t(v1, v2, v3, v4)
	}
}

// Partial
// partial func1 to func
func (t DelegateFunc1[TIn, TOut]) Partial(val TIn) DelegateFunc[TOut] {
	return func() TOut {
		return t(val)
	}
}

// Partial
// partial func2 to func1
func (t DelegateFunc2[TIn1, TIn2, TOut]) Partial(val TIn1) DelegateFunc1[TIn2, TOut] {
	return func(v2 TIn2) TOut {
		return t(val, v2)
	}
}
