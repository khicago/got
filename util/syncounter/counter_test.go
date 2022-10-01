package syncounter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter(t *testing.T) {
	c := MakeCounter(0)

	v, e1 := c.CountOne(nil)
	assert.Nil(t, e1, "count one failed, err should be NIL.")
	assert.Equal(t, int64(0), v)

	from, to, e2 := c.Count(100, nil)
	assert.Nil(t, e2, "count failed, err should be NIL.")
	assert.Equal(t, int64(1), from)
	assert.Equal(t, int64(100), to)

	from1, to1, e3 := c.Count(100, func() bool { return true })
	assert.Nil(t, e3, "count failed, err should be NIL.")
	assert.Equal(t, int64(0), from1)
	assert.Equal(t, int64(99), to1)
}

func TestCounterWithMaxVal(t *testing.T) {
	c := MakeCounter(10)

	v, e1 := c.CountOne(nil)
	assert.Nil(t, e1, "count one failed, err should be NIL.")
	assert.Equal(t, int64(0), v)

	_, _, e2 := c.Count(10, nil)
	assert.ErrorIs(t, e2, ErrCounterHasBeenExhausted, "err should be ErrCounterHasBeenExhausted.")

	from1, to1, e3 := c.Count(10, func() bool { return true })
	assert.Nil(t, e3, "count failed, err should be NIL.")
	assert.Equal(t, int64(0), from1)
	assert.Equal(t, int64(9), to1)
}
