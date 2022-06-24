package typer

type (
	Predicate[TVal any] func(val TVal) bool

	Action[T any]            func() T
	Func1[TIn any, TOut any] func(TIn) TOut
)

func Partial[TIn, TOut any](t Func1[TIn, TOut], val TIn) Action[TOut] {
	return func() TOut {
		return t(val)
	}
}
