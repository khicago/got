package proretry

import (
	"time"

	"github.com/khicago/got/util/delegate"
)

// Backoff
// - 退避算法
type Backoff delegate.Convert[time.Duration, time.Duration]

// ConstantBackoff
// - 按照指定的重试间隔进行退避。
func ConstantBackoff(interval time.Duration) Backoff {
	return func(time.Duration) time.Duration {
		return interval
	}
}

// LinearBackoff
// - 按照指定的初始重试间隔进行线性退避。
func LinearBackoff(initInterval time.Duration) Backoff {
	return func(prevInterval time.Duration) time.Duration {
		return initInterval + prevInterval
	}
}

// ExponentialBackoff
// - 按照指定的初始重试间隔进行指数退避
func ExponentialBackoff(initInterval time.Duration) Backoff {
	return func(prevInterval time.Duration) time.Duration {
		if prevInterval == 0 {
			return initInterval
		}
		return prevInterval * 2
	}
}

// FibonacciBackoff
// - 按照指定的初始重试间隔进行斐波那契退避
func FibonacciBackoff(initInterval time.Duration) Backoff {
	var prev int64 = 0
	var current int64 = 1
	return func(prevInterval time.Duration) time.Duration {
		prev, current = current, prev+current
		return initInterval * time.Duration(current)
	}
}
