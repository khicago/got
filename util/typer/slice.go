package typer

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/khicago/got/util/delegate"
	"golang.org/x/exp/constraints"

	"github.com/bytedance/gopkg/lang/fastrand"
)

func SliceForeach[TSliceVal any](slice []TSliceVal, foreachFn delegate.Action1[TSliceVal]) {
	for _, val := range slice {
		foreachFn.TryCall(val)
	}
}

func SliceForeachI[TSliceVal any](slice []TSliceVal, foreachFn delegate.Action2[TSliceVal, int]) {
	for i, val := range slice {
		foreachFn.TryCall(val, i)
	}
}

func SliceFirst[TVal comparable](slice []TVal, val TVal) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func SliceFirstMatch[TVal comparable](slice []TVal, pred delegate.Predicate[TVal]) int {
	for i, v := range slice {
		if pred(v) {
			return i
		}
	}
	return -1
}

func SliceFilter[TVal comparable](slice []TVal, pred delegate.Predicate[TVal]) []TVal {
	ret := make([]TVal, 0, len(slice)/2)
	for i := range slice {
		v := slice[i]
		if pred(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func SliceRand[TVal comparable](slice []TVal, defaultVal TVal) TVal {
	n := len(slice)
	if n == 0 {
		return defaultVal
	}
	i := fastrand.Intn(n)
	return slice[i]
}

func SliceContains[TVal comparable](slice []TVal, val TVal) bool {
	return SliceFirst(slice, val) >= 0
}

func SliceMap[TFrom, TTo any](from []TFrom, mapFn delegate.Convert[TFrom, TTo]) []TTo {
	ret := make([]TTo, 0, len(from))
	for _, valFrom := range from {
		valTo := mapFn(valFrom)
		ret = append(ret, valTo)
	}
	return ret
}

func SliceReduce[TSliceVal, TTarget any](slice []TSliceVal, reduceFn func(TSliceVal, TTarget) TTarget, defaultVal TTarget) TTarget {
	target := defaultVal
	for _, valFrom := range slice {
		target = reduceFn(valFrom, target)
	}
	return target
}

func SliceLast[TSliceVal any](slice []TSliceVal) TSliceVal {
	return slice[len(slice)-1]
}

func SliceTryGet[TSliceVal any](slice []TSliceVal, i int, defaultVal TSliceVal) TSliceVal {
	if i < 0 || i >= len(slice) {
		return defaultVal
	}
	return slice[i]
}

func SlicePadRight[TSliceVal any](slice []TSliceVal, length int, padVal TSliceVal) []TSliceVal {
	for len(slice) < length {
		slice = append(slice, padVal)
	}
	return slice
}

func SliceSort[TSliceVal constraints.Ordered](slice []TSliceVal) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
}

func SliceReverse[TSliceVal any](data []TSliceVal) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func SliceMax[TSliceVal constraints.Ordered](data []TSliceVal) (ret TSliceVal) {
	if len(data) == 0 {
		return
	}
	ret = data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > ret {
			ret = data[i]
		}
	}
	return ret
}

func SliceMin[TSliceVal constraints.Ordered](data []TSliceVal) (ret TSliceVal) {
	if len(data) == 0 {
		return
	}
	ret = data[0]
	for i := 1; i < len(data); i++ {
		if data[i] < ret {
			ret = data[i]
		}
	}
	return ret
}

func SliceToTrueMap[TSliceVal comparable](list []TSliceVal) map[TSliceVal]struct{} {
	m := make(map[TSliceVal]struct{}, len(list))
	for _, v := range list {
		m[v] = struct{}{}
	}
	return m
}

func SliceDiff[TSliceVal comparable](oldLst, newLst []TSliceVal) (toAdd, toRemove []TSliceVal) {
	// Threshold to decide between loop and map approach
	threshold := 10

	// Determine the approach based on the threshold
	var oldMap, newMap map[TSliceVal]struct{} = nil, nil
	if len(oldLst) > threshold {
		oldMap = SliceToTrueMap(oldLst)
	}

	if len(newLst) > threshold {
		newMap = SliceToTrueMap(newLst)
	}

	// Find elements to add
	for _, v := range newLst {
		if oldMap != nil {
			if _, exists := oldMap[v]; !exists {
				toAdd = append(toAdd, v)
			}
		} else if !SliceContains(oldLst, v) {
			toAdd = append(toAdd, v)
		}
	}
	// Find elements to remove
	for _, v := range oldLst {
		if newMap != nil {
			if _, exists := newMap[v]; !exists {
				toRemove = append(toRemove, v)
			}
		} else if !SliceContains(newLst, v) {
			toRemove = append(toRemove, v)
		}
	}
	return toAdd, toRemove
}

// IsSlice checks if the given variable is a slice.
func IsSlice(v any) bool {
	// Use reflect.ValueOf to get the reflection value of the variable
	val, ok := v.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(v)
	}
	// Check if the kind of the value is reflect.Slice
	return val.Kind() == reflect.Slice || val.Kind() == reflect.Array
}

// Is2DSlice checks if the given variable is a 2D array.
func Is2DSlice(v any) bool {
	// Use reflect.ValueOf to get the reflection value of the variable
	val, ok := v.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(v)
	}

	// Check if the kind of the value is reflect.Array or reflect.Slice
	if !IsSlice(val) {
		return false
	}

	// Check if the element type is also an array
	elemType := val.Type().Elem()
	return elemType.Kind() == reflect.Array || elemType.Kind() == reflect.Slice
}

// Flatten2DSliceGeneric flattens a 2D array of type T into a 1D array.
func Flatten2DSliceGeneric[T any](input ...[]T) []T {
	flattened := make([]T, 0)
	for _, subArray := range input {
		flattened = append(flattened, subArray...)
	}
	return flattened
}

// Flatten2DSlice flattens a 2D array into a 1D array.
func Flatten2DSlice(input any) ([]any, error) {
	val := reflect.ValueOf(input)
	if !IsSlice(val) {
		return nil, fmt.Errorf("input is not an array or slice")
	}

	elemType := val.Type().Elem()
	if elemType.Kind() != reflect.Array && elemType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a 2D slice")
	}

	flattened := make([]any, 0)
	for i := 0; i < val.Len(); i++ {
		subArray := val.Index(i)
		for j := 0; j < subArray.Len(); j++ {
			flattened = append(flattened, subArray.Index(j).Interface())
		}
	}

	return flattened, nil
}

// FlattenNestedSlices flattens nested slices or arrays into a single-dimensional slice.
//
// It takes two parameters:
//   - input: the input data of any type, which may contain nested slices or arrays.
//   - depth: the target flattening depth, controlling the number of levels to flatten.
//
// If the target depth is 0, no flattening is performed, and the input is returned as a slice.
//
// The function recursively flattens the input data up to the specified depth. If the current
// depth reaches the target depth, the elements are added to the result slice without further
// flattening.
//
// It returns a []any slice containing the flattened elements.
func FlattenNestedSlices(input any, depth int) []any {
	if depth <= 0 {
		return []any{input}
	}

	val := reflect.ValueOf(input)
	if !IsSlice(val) {
		return []any{input}
	}

	result := make([]any, 0)

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)

		if item.Kind() == reflect.Interface {
			actualItem := item.Elem()

			if actualItem.Kind() == reflect.Array || actualItem.Kind() == reflect.Slice {
				flattenedItem := FlattenNestedSlices(actualItem.Interface(), depth-1)
				result = append(result, flattenedItem...)
			} else {
				result = append(result, actualItem.Interface())
			}
		} else if item.Kind() == reflect.Array || item.Kind() == reflect.Slice {
			flattenedItem := FlattenNestedSlices(item.Interface(), depth-1)
			result = append(result, flattenedItem...)
		} else {
			result = append(result, item.Interface())
		}
	}

	return result
}

func DoAsSlice[T any](input any, cb func(val T) error) error {
	val := reflect.ValueOf(input)
	if !IsSlice(val) {
		if v, ok := input.(T); ok {
			return cb(v)
		}
		return fmt.Errorf("input type mismatch: expected %T, got %T", *new(T), input)
	}
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if v, ok := elem.(T); ok {
			if err := cb(v); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("element type mismatch: expected %T, got %T", *new(T), elem)
		}
	}
	return nil
}

// no needs to provide stack fn
//
//func SlicePushTail[TSliceVal any](slicePtr *[]TSliceVal, val TSliceVal) {
//	*slicePtr = append(*slicePtr, val)
//}
//
//func SlicePopTail[TSliceVal any](slicePtr *[]TSliceVal) (TSliceVal, error) {
//	l := len(*slicePtr)
//	if l == 0 {
//		return ZeroVal[TSliceVal](), errors.New("insufficient stack length")
//	}
//	v := (*slicePtr)[l-1]
//	*slicePtr = (*slicePtr)[:l-1]
//	return v, nil
//}
