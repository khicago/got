package typer

import "unsafe"

func AssertType[TAssert, TVal any](v TVal) bool {
	_, ok := any(v).(TAssert)
	return ok
}

func AssertNil(v any) bool {
	return (*struct {
		p    uintptr
		data unsafe.Pointer
	})(unsafe.Pointer(&v)).data == nil
}

func AssertNotNil(v any) bool {
	return !AssertNil(v)
}

func AssertZeroVal[T comparable](v T) bool {
	return v == ZeroVal[T]()
}
