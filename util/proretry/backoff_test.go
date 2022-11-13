package proretry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstantBackoff(t *testing.T) {
	backoff := ConstantBackoff(time.Second)
	assert.Equal(t, time.Second, backoff(0))
	assert.Equal(t, time.Second, backoff(time.Second))
	assert.Equal(t, time.Second, backoff(2*time.Second))
}

func TestLinearBackoff(t *testing.T) {
	backoff := LinearBackoff(time.Second)
	assert.Equal(t, time.Second, backoff(0))
	assert.Equal(t, 2*time.Second, backoff(time.Second))
	assert.Equal(t, 3*time.Second, backoff(2*time.Second))
}

func TestExponentialBackoff(t *testing.T) {
	backoff := ExponentialBackoff(time.Second)
	assert.Equal(t, time.Second, backoff(0))
	assert.Equal(t, 2*time.Second, backoff(time.Second))
	assert.Equal(t, 4*time.Second, backoff(2*time.Second))
	assert.Equal(t, 8*time.Second, backoff(4*time.Second))
}

func TestFibonacciBackoff(t *testing.T) {
	backoff := FibonacciBackoff(time.Second)
	assert.Equal(t, time.Second, backoff(0))
	assert.Equal(t, 2*time.Second, backoff(time.Second))
	assert.Equal(t, 3*time.Second, backoff(2*time.Second))
	assert.Equal(t, 5*time.Second, backoff(3*time.Second))
	assert.Equal(t, 8*time.Second, backoff(4*time.Second))
}
