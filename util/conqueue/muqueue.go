package conqueue

import (
	"sync"
)

type MuRing struct {
	sz   uint32
	dLen uint32
	cur  uint32
	head uint32
	mu   sync.Mutex
	arr  []interface{}
}

func newMuRing(size uint32) ConcurrentQueue {
	if size < MinRingSize32 || size > MaxRingSize32 {
		panic(ErrRingSizeOutOfRange)
	}
	sz := size
	dLen := sz + 1
	return &MuRing{
		mu:   sync.Mutex{},
		arr:  make([]interface{}, dLen),
		sz:   sz,
		dLen: dLen,
	}
}

func (cr *MuRing) Cap() int {
	return int(cr.sz)
}

func (cr *MuRing) Len() int {
	return int((cr.head + cr.dLen - cr.cur) % cr.dLen)
}

func (cr *MuRing) Empty() bool {
	return cr.cur == cr.head
}

func (cr *MuRing) Full() bool {
	return cr.forward(cr.head) == cr.cur
}

func (cr *MuRing) Push(c interface{}) error {
	cr.mu.Lock()
	if cr.Full() {
		cr.mu.Unlock()
		return ErrFull
	}
	cr.arr[cr.head] = c
	cr.head = cr.forward(cr.head)
	cr.mu.Unlock()
	return nil
}

func (cr *MuRing) Pop() (item interface{}, err error) {
	cr.mu.Lock()
	if cr.Empty() {
		cr.mu.Unlock()
		return nil, ErrEmpty
	}
	c := cr.arr[cr.cur]
	cr.arr[cr.cur] = nil
	cr.cur = cr.forward(cr.cur)
	cr.mu.Unlock()
	return c, nil
}

func (cr *MuRing) forward(pos uint32) (to uint32) {
	return (pos + 1) % cr.dLen
}
