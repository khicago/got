package conqueue

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestRingPushAndPop(t *testing.T) {
	queue := NewRing64(uint64(3))

	testGroup := []interface{}{
		1,
		' ',
		"1",
		"1234 aaa",
		struct{}{},
	}

	for _, c := range testGroup {
		errPush := queue.Push(c)
		a, errPop := queue.Pop()
		assert.Nil(t, errPush)
		assert.Nil(t, errPop)
		assert.Equal(t, c, a)
	}
	assert.Equal(t, 0, queue.(*Ring64).Len())
}

func TestRingPushAndPopConcurrentPus(t *testing.T) {
	size := 100000
	co := 5
	lineSize := size / co
	queue := NewRing64(uint64(size))

	wg := sync.WaitGroup{}
	wg.Add(co)
	te := func(from, to int) {
		defer wg.Done()
		for ; from < to; from++ {
			err := queue.Push(from)
			if !assert.Nil(t, err) {
				return
			}
		}
	}

	for i := 0; i < co; i++ {
		go te(co*lineSize, (1+co)*lineSize)
	}
	wg.Wait()
	assert.Equal(t, size, queue.(*Ring64).Len())
}

func TestRingPushAndPopConcurrentPushPop(t *testing.T) {
	size := 100000
	co := 5
	lineSize := size / co
	queue := NewRing64(uint64(size))

	wg := sync.WaitGroup{}
	wg.Add(co)
	te := func(from, to int) {
		defer wg.Done()
		for ; from < to; from++ {
			err := queue.Push(from)
			if !assert.Nil(t, err) {
				return
			}
			_, err = queue.Pop()
			if !assert.Nil(t, err) {
				return
			}
		}
	}

	for i := 0; i < co; i++ {
		go te(co*lineSize, (1+co)*lineSize)
	}
	wg.Wait()
	assert.Equal(t, 0, queue.(*Ring64).Len())

	wg = sync.WaitGroup{}
	wg.Add(co)
	te2 := func(from, to int) {
		defer wg.Done()
		for ; from < to; from++ {
			err := queue.Push(from)
			if !assert.Nil(t, err) {
				return
			}
			if from%2 == 0 {
				_, err = queue.Pop()
				if !assert.Nil(t, err) {
					return
				}
			}
		}
	}

	for i := 0; i < co; i++ {
		go te2(co*lineSize, (1+co)*lineSize)
	}
	wg.Wait()
	assert.Equal(t, 100000/2, queue.(*Ring64).Len())
}

func BenchmarkRingPushAndPop(b *testing.B) {
	size := uint64(b.N / 5)
	if size == 0 {
		size = 1
	}
	queue := NewRing64(size)
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
	assert.Equal(b, 0, queue.(*Ring64).Len())
}

func BenchmarkRingPush(b *testing.B) {
	size := uint64(b.N)
	if size == 0 {
		size = 1
	}
	queue := NewRing64(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := queue.Push(123456)
			if !assert.Nil(b, err) {
				return
			}
		}
	})
	assert.Equal(b, b.N, queue.(*Ring64).Len())
}
