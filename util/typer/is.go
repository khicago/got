package typer

import (
	"reflect"
	"unsafe"
)

func IsType[TAssert, TVal any](v TVal) bool {
	_, ok := any(v).(TAssert)
	return ok
}

func IsNil(v any) bool {
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

func IsNotNil(v any) bool {
	return !IsNil(v)
}

func IsZero[T comparable](v T) bool {
	return ZeroVal[T]() == v
}
