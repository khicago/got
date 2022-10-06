package conqueue

import (
	"fmt"
	"math"
	"runtime"
	"sync/atomic"
)

type Ring64 struct {
	cur   uint64
	_     [56]byte
	head  uint64
	_     [56]byte
	cap   uint64
	dLen  uint64
	bound uint64
	d     []interface{}
}

// 支持 64 位 size
func NewRing64(size uint64) ConcurrentQueue {
	if size < MinRingSize64 || size > MaxRingSize64 {
		panic(ErrRingSizeOutOfRange)
	}
	sz := size
	dLen := sz + 1
	return &Ring64{
		d:     make([]interface{}, dLen),
		cap:   sz,
		dLen:  dLen,
		bound: (math.MaxUint32 / dLen) * dLen,
	}
}

func (r *Ring64) Cap() int {
	return int(r.cap)
}

func (r *Ring64) Len() int {
	return int((r.head + r.dLen - r.cur) % r.dLen)
}

func (r *Ring64) Empty() bool {
	return r.cur == r.head
}

func (r *Ring64) Full() bool {
	return r.forward(r.head)%r.dLen == r.cur
}

func (r *Ring64) Push(c interface{}) error {
	for {
		cur, head := r.cur, r.head
		newHead := r.forward(head)
		if r.ind(newHead) == r.ind(cur) {
			fmt.Println("to%r.dLen == r.cur", newHead, r.dLen, r.cur)
			return ErrFull
		}
		if atomic.CompareAndSwapUint64(&r.head, head, newHead) {
			r.d[r.ind(head)] = c
			return nil
		}
		runtime.Gosched()
	}
}

func (r *Ring64) Pop() (item interface{}, err error) {
	for {
		cur := r.cur
		if cur == r.head {
			return nil, ErrEmpty
		}
		newCur := r.forward(cur)
		if atomic.CompareAndSwapUint64(&r.cur, cur, newCur) {
			return r.d[r.ind(cur)], nil
		}
		runtime.Gosched()
	}
}

func (r *Ring64) forward(pos uint64) (to uint64) {
	return (pos + 1) % r.bound
}

func (r *Ring64) ind(pos uint64) uint64 {
	return pos % r.dLen
}
