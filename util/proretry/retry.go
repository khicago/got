package proretry

import (
	"errors"
	"time"

	"github.com/khicago/got/util/delegate"
)

type (
	RetryFunc delegate.Handler

	// RetryOption
	// - Run 函数的选项
	RetryOption func(*options)

	options struct {
		initInterval  time.Duration
		retryableErrs delegate.Func1[error, bool]
	}
)

func WithInitInterval(interval time.Duration) RetryOption {
	return func(o *options) {
		o.initInterval = interval
	}
}

// WithRetryableErrs
// 设置 Run 函数的可重试错误列表
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

// WithRetryableErrFunc
// 设置 Run 函数的可重试错误判定函数
func WithRetryableErrFunc(f func(error) bool) RetryOption {
	return func(o *options) {
		o.retryableErrs = f
	}
}

// Run 执行指定的函数，并按照指定的退避算法进行重试。
//
// 参数：
//
//	fn: 待执行的函数。
//	maxRetries: 最大重试次数。
//	backoff: 退避算法。
//
// 返回值：
//
//	执行函数的错误，如果执行成功则返回 nil。
func Run(fn RetryFunc, maxRetries int, backoff Backoff, opts ...RetryOption) error {
	// 设置默认值
	option := &options{
		initInterval: 80 * time.Millisecond,
		retryableErrs: func(error) bool {
			return true
		},
	}
	// 应用选项
	for _, o := range opts {
		o(option)
	}

	var err error
	interval := option.initInterval
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if !option.retryableErrs.TryCall(err, true) {
			return err
		}
		time.Sleep(interval)
		interval = backoff(interval)
	}
	return err
}
