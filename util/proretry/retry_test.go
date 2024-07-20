package proretry

import (
	"errors"
	"fmt"
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
			maxRetries: 3,
			backoff:    ConstantBackoff(time.Millisecond * 10),
			expectedErr: &RetryError{
				Attempts: 3,
				LastErr:  errors.New("error"),
			},
			expectedCalls: 3,
		},
		{
			name: "retryable error",
			fn: func() error {
				return retryableErr
			},
			maxRetries: 3,
			backoff:    ConstantBackoff(time.Millisecond * 10),
			options:    []RetryOption{WithRetryableErrs(retryableErr)},
			expectedErr: &RetryError{
				Attempts: 3,
				LastErr:  retryableErr,
			},
			expectedCalls: 3,
		},
		{
			name: "non-retryable error",
			fn: func() error {
				return errors.New("error")
			},
			maxRetries: 3,
			options: []RetryOption{
				WithRetryableErrs(errors.New("other error")),
				WithBackoff(ConstantBackoff(time.Millisecond * 10)),
			},
			expectedErr: &RetryError{
				Attempts: 0,
				LastErr:  errors.New("error"),
			},
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
			err := Run(test.fn, test.maxRetries, test.options...)
			assert.Equal(t, fmt.Sprintf("%v", test.expectedErr), fmt.Sprintf("%v", err))
			assert.Equal(t, test.expectedCalls, calls)
		})
	}
}
