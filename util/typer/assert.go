package typer

import (
	"reflect"
	"unsafe"
)

func AssertType[TAssert, TVal any](v TVal) bool {
	_, ok := any(v).(TAssert)
	return ok
}

func AssertNil(v any) bool {
	if (*struct {
		p    uintptr
		data unsafe.Pointer
	})(unsafe.Pointer(&v)).data == nil {
		return true
	}

	valueOf := reflect.ValueOf(v)
	k := valueOf.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return valueOf.IsNil()
	default:
		return v == nil
	}
}

func AssertNotNil(v any) bool {
	return !AssertNil(v)
}

func AssertZeroVal[T comparable](v T) bool {
	return v == ZeroVal[T]()
}
