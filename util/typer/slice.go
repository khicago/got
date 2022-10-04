package typer

import (
	"github.com/khicago/got/util/delegate"
	"golang.org/x/exp/constraints"
	"sort"
)

func SliceFirst[TVal comparable](slice []TVal, val TVal) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func SliceContains[TVal comparable](slice []TVal, val TVal) bool {
	return SliceFirst(slice, val) >= 0
}

func SliceMap[TFrom, TTo any](from []TFrom, mapFn delegate.Map[TFrom, TTo]) []TTo {
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
