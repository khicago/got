package typer

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

func SliceMap[TFrom, TTo any](from []TFrom, mapFn DelegateMap[TFrom, TTo]) []TTo {
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

func SlicePadRight[TSliceVal any](slice []TSliceVal, length int, padVal TSliceVal) []TSliceVal {
	for len(slice) < length {
		slice = append(slice, padVal)
	}
	return slice
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
