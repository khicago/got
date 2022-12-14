package reflecto

import (
	"fmt"
	"reflect"
)

type SlicePtrReflector struct {
	ptr        any
	sliceValue *reflect.Value

	elemType reflect.Type
}

func (ptrRef *SlicePtrReflector) Ptr() int {
	return ptrRef.sliceValue.Len()
}

func (ptrRef *SlicePtrReflector) Len() int {
	return ptrRef.sliceValue.Len()
}

func (ptrRef *SlicePtrReflector) ItemType() reflect.Type {
	if ptrRef.elemType != nil {
		return ptrRef.elemType
	}
	ptrRef.elemType = ptrRef.sliceValue.Type().Elem()
	return ptrRef.elemType
}

func (ptrRef *SlicePtrReflector) Read(i int, outPtr any) error {
	vOut := reflect.ValueOf(outPtr)
	if vOut.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid outter, should be a pointer of elem %v", vOut.Kind())
	}
	itemVal := ptrRef.sliceValue.Index(i)

	if vOut.Type().Elem() == ptrRef.ItemType() {
		vOut.Elem().Set(itemVal)
	} else if vOut.Type() == ptrRef.ItemType() {
		vOut.Elem().Set(itemVal.Elem())
	} else {
		return fmt.Errorf("invalid outter, type not match, expect %v got %v", ptrRef.ItemType(), vOut.Type().Elem())
	}

	return nil
}

func NewSlicePtrReflector(slicePtr any) (*SlicePtrReflector, error) {
	value, err := ToWriteableSliceValue(slicePtr)
	if err != nil {
		return nil, err
	}
	return &SlicePtrReflector{
		ptr:        slicePtr,
		sliceValue: value,
	}, nil
}

func ToWriteableSliceValue(slicePointer any) (*reflect.Value, error) {
	vSlicePtr := reflect.ValueOf(slicePointer)
	// to make sure it are addressable
	if vSlicePtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", vSlicePtr.Type())
	}
	vSlice := vSlicePtr.Elem()
	if vSlice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", vSlice.Type())
	}
	return &vSlice, nil
}

func GetSliceElementType(slice any) (reflect.Type, error) {
	ty := reflect.TypeOf(slice)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	if ty.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", ty)
	}
	return ty.Elem(), nil
}

func SliceContains(slice any, target any) bool {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < val.Len(); i++ {
		if target == val.Index(i).Interface() {
			return true
		}
	}
	return false
}
