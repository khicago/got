package typer

import (
	"github.com/khicago/got/util/delegate"
	"golang.org/x/exp/constraints"
)

func Keys[TKey comparable, TVal any](m map[TKey]TVal) []TKey {
	keys := make([]TKey, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func KeysSorted[TKey constraints.Ordered, TVal any](m map[TKey]TVal) []TKey {
	keys := Keys(m)
	SliceSort(keys)
	return keys
}

func MapForEachOrderly[TKey constraints.Ordered, TVal any](m map[TKey]TVal, traver func(key TKey, val TVal)) {
	keys := KeysSorted(m)
	for _, key := range keys {
		traver(key, m[key])
	}
}

func MapMap[TKey comparable, TVal any, TKey1 comparable, TVal1 any](m map[TKey]TVal, keyConv delegate.Convert[TKey, TKey1], valConv delegate.Convert[TVal, TVal1]) map[TKey1]TVal1 {
	ret := make(map[TKey1]TVal1)
	for key := range m {
		ret[keyConv(key)] = valConv(m[key])
	}
	return ret
}

func MapDump[TKey comparable, TVal any](m map[TKey]TVal, dumper delegate.Func2[TKey, TVal, string]) []string {
	ret := make([]string, 0, len(m))
	for key, value := range m {
		ret = append(ret, dumper(key, value))
	}
	return ret
}
