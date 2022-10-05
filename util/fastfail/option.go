package fastfail

import (
	"errors"
	"time"
)

type (
	conf struct {
		Retry        int
		WaitInit     time.Duration
		WaitIncrease time.Duration
	}
)

var (
	ErrLockFailed   = errors.New("fast-fail: lock failed")
	ErrUnlockFailed = errors.New("fast-fail: unlock failed")
)

const (
	RetryDefault  int = 3
	RetryUntilGet int = -1
)

// DefaultFastFailConf is the set of default configs of fastFailLock
// It can be set to change the default behavior of all FastFailLock
var DefaultFastFailConf = conf{
	Retry:        RetryDefault,
	WaitInit:     time.Microsecond * time.Duration(100),
	WaitIncrease: time.Microsecond * time.Duration(201),
}

func (c conf) Pipe(opts ...FastFailLockOption) (ret conf) {
	ret = c
	for _, opt := range opts {
		ret = opt(ret)
	}
	return ret
}

// FFOptRetry sets the retry count of the FastFailLock
// If the given retryCount <= 0, the lock will spin forever
func FFOptRetry(retryCount int) FastFailLockOption {
	return func(org conf) conf {
		org.Retry = retryCount
		return org
	}
}

func FFOptWaitInit(waitInit time.Duration) FastFailLockOption {
	return func(org conf) conf {
		org.WaitInit = waitInit
		return org
	}
}

func FFOptWaitIncrease(waitIncrease time.Duration) FastFailLockOption {
	return func(org conf) conf {
		org.WaitIncrease = waitIncrease
		return org
	}
}
