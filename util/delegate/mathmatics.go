package delegate

type (
	Compare[TVal any] func(TVal, TVal) int
)
