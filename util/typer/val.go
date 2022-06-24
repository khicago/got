package typer

func ZeroVal[T any]() (v T) {
	return
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
