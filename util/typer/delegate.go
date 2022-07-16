package typer

type (
	DelegateAction[T any]          func() T
	DelegateMap[TIn any, TOut any] func(TIn) TOut

	Predicate[TVal any] DelegateMap[TVal, bool]
)

func Partial[TIn, TOut any](t DelegateMap[TIn, TOut], val TIn) DelegateAction[TOut] {
	return func() TOut {
		return t(val)
	}
}
