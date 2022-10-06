package conqueue

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMuRingPushAndPop(t *testing.T) {
	queue := newMuRing(uint32(3))

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
	assert.Equal(t, 0, queue.(*MuRing).Len())
}

func TestMuRingPushAndPopConcurrentPus(t *testing.T) {
	size := 100000
	co := 5
	lineSize := size / co
	queue := newMuRing(uint32(size))

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
	assert.Equal(t, size, queue.(*MuRing).Len())
}

func TestMuRingPushAndPopConcurrentPushPop(t *testing.T) {
	size := 100000
	co := 5
	lineSize := size / co
	queue := newMuRing(uint32(size))

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
	assert.Equal(t, 0, queue.(*MuRing).Len())

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
	assert.Equal(t, 100000/2, queue.(*MuRing).Len())
}

func BenchmarkMuRingPushAndPop(b *testing.B) {
	size := uint32(b.N / 5)
	if size == 0 {
		size = 1
	}
	queue := newMuRing(size)
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
	assert.Equal(b, 0, queue.(*MuRing).Len())
}


func BenchmarkMuRingPush(b *testing.B) {
	size := uint32(b.N)
	if size == 0 {
		size = 1
	}
	queue := newMuRing(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := queue.Push(123456)
			if !assert.Nil(b, err) {
				return
			}
		}
	})
	assert.Equal(b, b.N, queue.(*MuRing).Len())
}
