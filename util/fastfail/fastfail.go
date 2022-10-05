package fastfail

import (
	"sync/atomic"
	"time"
)

type (
	FastFailLock struct {
		int32
	}

	FastFailLockOption func(org conf) conf
)

func (ffl *FastFailLock) Lock(opts ...FastFailLockOption) error {
	c := DefaultFastFailConf.Pipe(opts...)

	if execCount := c.Retry + 1; execCount > 0 {
		for i := 1; i <= execCount; i++ {
			if atomic.CompareAndSwapInt32(&ffl.int32, 0, 1) {
				return nil
			}
			time.Sleep(c.WaitInit + c.WaitIncrease*time.Duration(i))
		}
		return ErrLockFailed
	}

	for {
		if atomic.CompareAndSwapInt32(&ffl.int32, 0, 1) {
			return nil
		}
		time.Sleep(c.WaitInit)
	}
}

func (ffl *FastFailLock) Unlock() {
	if !atomic.CompareAndSwapInt32(&ffl.int32, 1, 0) {
		panic(ErrUnlockFailed)
	}
}

func NewFastFailLock() *FastFailLock {
	return &FastFailLock{}
}
