package conqueue

import "errors"

type (
	ConcurrentQueue interface {
		Cap() int
		Len() int
		Empty() bool
		Full() bool
		Push(c interface{}) error
		Pop() (ret interface{}, err error)
	}
)

var (
	ErrFull  = errors.New("queue is full")
	ErrEmpty = errors.New("queue is empty")

	ErrRingSizeOutOfRange = errors.New("ring size is out of range")
)

const (
	MinRingSize32 = 1
	MaxRingSize32 = 1 << 30
	MinRingSize64 = 1
	MaxRingSize64 = 1 << 60
)
