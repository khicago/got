package proretry

import (
	"errors"
	"fmt"
	"time"

	"github.com/khicago/got/util/delegate"
)

type (
	// RetryFunc 是一个函数类型，用于定义需要重试的函数。
	RetryFunc delegate.Handler

	// RetryOption 是一个函数类型，用于配置 Run 函数的选项。
	RetryOption func(*options)

	// options 包含 Run 函数的配置选项。
	options struct {
		initInterval  time.Duration
		retryableErrs delegate.Func1[error, bool]
		backoff       Backoff
	}

	// RetryError 包含重试次数和最后一次错误信息。
	RetryError struct {
		Attempts int
		LastErr  error
	}
)

// DefaultRetryInterval 是默认的重试间隔时间。
var DefaultRetryInterval = 80 * time.Millisecond

// Error 实现 error 接口，返回包含重试次数和最后一次错误信息的字符串。
func (r *RetryError) Error() string {
	return fmt.Sprintf("after %d attempts, last error: %v", r.Attempts, r.LastErr)
}

// LastError 返回最后一次错误信息。
func (r *RetryError) LastError() error {
	return r.LastErr
}

// WithInitInterval 设置初始重试间隔时间。
func WithInitInterval(interval time.Duration) RetryOption {
	return func(o *options) {
		o.initInterval = interval
	}
}

// WithRetryableErrs 设置可重试的错误列表。
func WithRetryableErrs(errs ...error) RetryOption {
	return func(o *options) {
		o.retryableErrs = func(err error) bool {
			for _, e := range errs {
				if errors.Is(err, e) {
					return true
				}
			}
			return false
		}
	}
}

// WithRetryableErrFunc 设置可重试错误的判定函数。
func WithRetryableErrFunc(f func(error) bool) RetryOption {
	return func(o *options) {
		o.retryableErrs = f
	}
}

// WithBackoff 设置退避算法。
func WithBackoff(backoff Backoff) RetryOption {
	return func(o *options) {
		o.backoff = backoff
	}
}

// Run 执行指定的函数，并按照指定的退避算法进行重试。
//
// 参数：
//   - fn: 待执行的函数。
//   - maxRetries: 最大重试次数。
//   - opts: 可选的配置选项。
//
// 返回值：
//   - 执行函数的错误，如果执行成功则返回 nil。
func Run(fn RetryFunc, maxRetries int, opts ...RetryOption) error {
	// 设置默认值
	option := &options{
		initInterval: DefaultRetryInterval,
		retryableErrs: func(error) bool {
			return true
		},
		backoff: FibonacciBackoff(DefaultRetryInterval),
	}
	// 应用选项
	for _, o := range opts {
		o(option)
	}

	interval := option.initInterval
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if !option.retryableErrs.TryCall(err, true) {
			return &RetryError{
				Attempts: i,
				LastErr:  err,
			}
		}
		time.Sleep(interval)
		interval = option.backoff(interval)
	}
	return &RetryError{
		Attempts: maxRetries,
		LastErr:  err,
	}
}
