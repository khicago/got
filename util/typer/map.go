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

func MapForEachOrderly[TKey constraints.Ordered, TVal any](m map[TKey]TVal, traver func(key TKey, val TVal)) {
	keys := Keys(m)
	SliceSort(keys)
	for _, key := range keys {
		traver(key, m[key])
	}
}

func MapDump[TKey comparable, TVal any](m map[TKey]TVal, dumper delegate.Func2[TKey, TVal, string]) []string {
	ret := make([]string, 0, len(m))
	for key, value := range m {
		ret = append(ret, dumper(key, value))
	}
	return ret
}
