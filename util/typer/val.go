package typer

import (
	"golang.org/x/exp/constraints"
)

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

func Or[T comparable](a, b T) T {
	if IsZero(a) {
		return b
	}
	return a
}

func Between[T constraints.Ordered](v T, begin, end T) bool {
	return v >= begin && v <= end
}
