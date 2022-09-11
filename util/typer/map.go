package typer

import (
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
