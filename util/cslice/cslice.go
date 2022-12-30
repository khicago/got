package cslice

import (
	"sync"
)

// CSlice is a concurrently safe slice for data type T
type CSlice[T any] struct {
	data []T
	mu   sync.Mutex
}

// New creates a new CSlice
func New[T any](cap int) CSlice[T] {
	return CSlice[T]{
		data: make([]T, 0, cap),
		mu:   sync.Mutex{},
	}
}

func (b *CSlice[T]) Add(value T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = append(b.data, value)
}

func (b *CSlice[T]) RemoveLast() *T {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.data) == 0 {
		return nil
	}

	item := b.data[len(b.data)-1]
	b.data = b.data[:len(b.data)-1]
	return &item
}

func (b *CSlice[T]) Length() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.data)
}

// Reset clears the concurrent_slice
func (b *CSlice[T]) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = nil
}

// FillTail adds multiple items to the concurrent_slice
func (b *CSlice[T]) FillTail(items []T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, item := range items {
		b.data = append(b.data, item)
	}
}
