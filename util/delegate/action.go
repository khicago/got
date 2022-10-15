package delegate

type (
	Action                      func()
	Action1[T any]              func(T)
	Action2[T1, T2 any]         func(T1, T2)
	Action3[T1, T2, T3 any]     func(T1, T2, T3)
	Action4[T1, T2, T3, T4 any] func(T1, T2, T3, T4)
)

// TryCall
// do nothing when action is nil
func (t Action) TryCall() {
	if t == nil {
		return
	}
	t()
}

// Partial
// partial func1 to func
func (t Action1[T]) Partial(val T) Action {
	return func() {
		t(val)
	}
}

// WrapAction wrap an action
func WrapAction[T any](fn Action1[T]) Action1[T] {
	return fn
}

// TryCall
// do nothing when action is nil
func (t Action1[T]) TryCall(val T) {
	if t == nil {
		return
	}
	t(val)
}

// Partial
// partial func1 to func
func (t Action2[T1, T2]) Partial(v1 T1) Action1[T2] {
	return func(v2 T2) {
		t(v1, v2)
	}
}

// TryCall
// do nothing when action is nil
func (t Action2[T1, T2]) TryCall(v1 T1, v2 T2) {
	if t == nil {
		return
	}
	t(v1, v2)
}

// Partial
// partial func1 to func
func (t Action3[T1, T2, T3]) Partial(v1 T1) Action2[T2, T3] {
	return func(v2 T2, v3 T3) {
		t(v1, v2, v3)
	}
}

// Partial
// partial func1 to func
func (t Action4[T1, T2, T3, T4]) Partial(v1 T1) Action3[T2, T3, T4] {
	return func(v2 T2, v3 T3, v4 T4) {
		t(v1, v2, v3, v4)
	}
}
