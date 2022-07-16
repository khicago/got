package typer

func ZeroVal[T any]() (v T) {
	return
}

func IsZero[T comparable](v T) bool {
	return ZeroVal[T]() == v
}

func Ptr[T any](v T) *T {
	return &v
}

func IfThen[T any](cond bool, a T, b T) T {
	if cond {
		return a
	}
	return b
}

func Or[T comparable](a, b T) T {
	if IsZero(a) {
		return b
	}
	return a
}

func Keys[TKey comparable, TVal any](m map[TKey]TVal) []TKey {
	keys := make([]TKey, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
