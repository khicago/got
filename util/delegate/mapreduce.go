package delegate

type (
	Convert[TIn, TOut any] Func1[TIn, TOut]
	Reduce[TIn, TOut any]  Func2[TIn, TOut, TOut]
)
