package conqueue

import (
	"fmt"
	"math"
	"runtime"
	"sync/atomic"
)

type Ring struct {
	cur   uint32
	_     [60]byte
	head  uint32
	_     [60]byte
	cap   uint32
	dLen  uint32
	bound uint32
	d     []interface{}
}

func NewRing(size uint32) ConcurrentQueue {
	if size < MinRingSize32 || size > MaxRingSize32 {
		panic(ErrRingSizeOutOfRange)
	}

	dLen := size + 1
	return &Ring{
		d:     make([]interface{}, dLen),
		cap:   size,
		dLen:  dLen,
		bound: (math.MaxUint32 / dLen) * dLen,
	}
}

func (r *Ring) Cap() int {
	return int(r.cap)
}

func (r *Ring) Len() int {
	return int((r.head + r.dLen - r.cur) % r.dLen)
}

func (r *Ring) Empty() bool {
	return r.cur == r.head
}

func (r *Ring) Full() bool {
	return r.forward(r.head)%r.dLen == r.cur
}

func (r *Ring) Push(c interface{}) error {
	for {
		cur, head := r.cur, r.head // todo: ?
		// 需要用较新的 head 和较旧的 cur 比较, 因此不可使用缓存 head (考虑编译优化)

		newHead := r.forward(head)
		if r.ind(newHead) == r.ind(cur) { // todo: to == cur ?
			fmt.Println("to%r.dLen == r.cur", newHead, r.dLen, r.cur)
			return ErrFull
		}
		if atomic.CompareAndSwapUint32(&r.head, head, newHead) {
			r.d[r.ind(head)] = c
			return nil
		}
		runtime.Gosched()
	}
}

func (r *Ring) Pop() (item interface{}, err error) {
	for {
		cur := r.cur
		if cur == r.head {
			return nil, ErrEmpty
		}
		newCur := r.forward(cur)
		if atomic.CompareAndSwapUint32(&r.cur, cur, newCur) {
			return r.d[r.ind(cur)], nil
		}
		runtime.Gosched()
	}
}

func (r *Ring) forward(pos uint32) (to uint32) {
	return (pos + 1) % r.bound
}

func (r *Ring) ind(pos uint32) uint32 {
	return pos % r.dLen
}
