package proretry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryDo(t *testing.T) {
	retryableErr := errors.New("given error")
	tests := []struct {
		name          string
		fn            RetryFunc
		maxRetries    int
		backoff       Backoff
		options       []RetryOption
		expectedErr   error
		expectedCalls int
	}{
		{
			name: "success",
			fn: func() error {
				return nil
			},
			maxRetries:    2,
			backoff:       ConstantBackoff(time.Millisecond * 10),
			expectedErr:   nil,
			expectedCalls: 1,
		},
		{
			name: "max retries exceeded",
			fn: func() error {
				return errors.New("error")
			},
			maxRetries:    3,
			backoff:       ConstantBackoff(time.Millisecond * 10),
			expectedErr:   errors.New("error"),
			expectedCalls: 3,
		},
		{
			name: "retryable error",
			fn: func() error {
				return retryableErr
			},
			maxRetries:    3,
			backoff:       ConstantBackoff(time.Millisecond * 10),
			options:       []RetryOption{WithRetryableErrs(retryableErr)},
			expectedErr:   retryableErr,
			expectedCalls: 3,
		},
		{
			name: "non-retryable error",
			fn: func() error {
				return errors.New("error")
			},
			maxRetries:    3,
			backoff:       ConstantBackoff(time.Millisecond * 10),
			options:       []RetryOption{WithRetryableErrs(errors.New("other error"))},
			expectedErr:   errors.New("error"),
			expectedCalls: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			fn := test.fn
			test.fn = func() error {
				calls++
				return fn()
			}
			err := Run(test.fn, test.maxRetries, test.backoff, test.options...)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedCalls, calls)
		})
	}
}
