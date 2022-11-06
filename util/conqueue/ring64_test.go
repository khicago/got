package conqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRing64PushAndPop(b *testing.B) {
	size := uint32(b.N / 5)
	if size == 0 {
		size = 1
	}
	queue := NewRing(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := queue.Push(123456)
			if !assert.Nil(b, err) {
				return
			}
			_, err = queue.Pop()
			if !assert.Nil(b, err) {
				return
			}
		}
	})
	assert.Equal(b, 0, queue.(*Ring).Len())
}

func BenchmarkRing64Push(b *testing.B) {
	size := uint32(b.N)
	if size == 0 {
		size = 1
	}
	queue := NewRing(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := queue.Push(123456)
			if !assert.Nil(b, err) {
				return
			}
		}
	})
	assert.Equal(b, b.N, queue.(*Ring).Len())
}
