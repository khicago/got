package syncounter

import (
	"errors"
	"github.com/khicago/got/util/delegate"
	"sync"
)

type Counter struct {
	mu      sync.Mutex
	counter int64
	max     int64
}

var (
	ErrCounterHasBeenExhausted = errors.New("counter has been exhausted")
)

// MakeCounter create a counter
//
// The upper limit of Counter determined by the max parameter,
// which is the maximum value of the top line int64 if set to
// a number less than or equal to 0.
func MakeCounter(max int64) Counter {
	return Counter{
		max: max,
	}
}

// CountOne return 1 positive int64 value
//
// when the condClear returns true, the counter will be reset.
func (c *Counter) CountOne(condClear delegate.Condition) (val int64, err error) {
	val, _, err = c.Count(1, condClear)
	return
}

// Count Get a set of consecutive positive int64 values and
// return their start and end values.
//
// When condClear returns true, the counter will be reset.
func (c *Counter) Count(offset int64, condClear delegate.Condition) (start, end int64, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if condClear != nil && condClear() {
		c.counter = 0
	}
	start = c.counter
	c.counter += offset
	if c.max > 0 && c.counter > c.max {
		return 0, 0, ErrCounterHasBeenExhausted
	}
	end = c.counter - 1
	return
}
