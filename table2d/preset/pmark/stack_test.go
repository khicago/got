package pmark_test

import (
	"testing"

	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/stretchr/testify/assert"
)

// TestStackPushPop - test pmark.Stack's Push and Pop methods
func TestStackPushPop(t *testing.T) {
	stack := pmark.Stack[int]{}
	// test empty len & peek & pop
	assert.Equal(t, 0, stack.Len())
	_, err := stack.Peek()
	assert.Error(t, err)
	_, err = stack.Pop()
	assert.Error(t, err)

	// push
	stack.Push(pmark.BucketsAngleL, 1)
	stack.Push(pmark.BucketsAngleR, 2)
	stack.Push(pmark.BucketsCurlyL, 3)
	stack.Push(pmark.BucketsCurlyR, 4)
	stack.Push(pmark.BucketsRoundL, 5)
	stack.Push(pmark.BucketsRoundR, 6)
	stack.Push(pmark.BucketsSquareL, 7)
	stack.Push(pmark.BucketsSquareR, 8)

	// test stack.Len()
	assert.Equal(t, 8, stack.Len())

	// test stack.PeekMark()
	m, err := stack.PeekMark()
	assert.NoError(t, err)
	assert.Equal(t, pmark.BucketsSquareR, m)

	// test stack.Peek()
	v, err := stack.Peek()
	assert.NoError(t, err)
	assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsSquareR, Val: 8}, v)

	// test stack.Pop()
	v, err = stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsSquareR, Val: 8}, v)
	assert.Equal(t, 7, stack.Len())
}

// TestStackConsume - test pmark.Stack.Consume() with testify
func TestStackConsume(t *testing.T) {
	// test stack.Consume()
	stack := pmark.Stack[int]{}
	err := stack.Consume(pmark.BucketsAngleL, 1, func(meet bool, v pmark.Pair[int]) {
		assert.False(t, meet)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsAngleL, Val: 1}, v.L)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.NIL}, v.R)
	})
	assert.NoError(t, err)
	err = stack.Consume(pmark.BucketsAngleR, 2, func(meet bool, v pmark.Pair[int]) {
		assert.True(t, meet)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsAngleL, Val: 1}, v.L)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsAngleR, Val: 2}, v.R)
	})
	assert.NoError(t, err)
	err = stack.Consume(pmark.BucketsCurlyL, 3, func(meet bool, v pmark.Pair[int]) {
		assert.False(t, meet)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsCurlyL, Val: 3}, v.L)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.NIL}, v.R)
	})
	assert.NoError(t, err)
	err = stack.Consume(pmark.BucketsRoundL, 4, func(meet bool, v pmark.Pair[int]) {
		assert.False(t, meet)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsRoundL, Val: 4}, v.L)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.NIL}, v.R)
	})
	assert.NoError(t, err)
	err = stack.Consume(pmark.BucketsRoundR, 5, func(meet bool, v pmark.Pair[int]) {
		assert.True(t, meet)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsRoundL, Val: 4}, v.L)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsRoundR, Val: 5}, v.R)
	})
	assert.NoError(t, err)

	err = stack.Consume(pmark.BucketsCurlyR, 6, func(meet bool, v pmark.Pair[int]) {
		assert.True(t, meet)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsCurlyL, Val: 3}, v.L)
		assert.Equal(t, pmark.Cell[int]{Mark: pmark.BucketsCurlyR, Val: 6}, v.R)
	})
	assert.NoError(t, err)

	// verify stack.Len()
	assert.Equal(t, 0, stack.Len())

	// verify stack.Results
	assert.Equal(t, 3, len(stack.Results))
	assert.Equal(t, pmark.Pair[int]{
		L: pmark.Cell[int]{Mark: pmark.BucketsAngleL, Val: 1},
		R: pmark.Cell[int]{Mark: pmark.BucketsAngleR, Val: 2},
	}, stack.Results[0])
	assert.Equal(t, pmark.Pair[int]{
		L: pmark.Cell[int]{Mark: pmark.BucketsRoundL, Val: 4},
		R: pmark.Cell[int]{Mark: pmark.BucketsRoundR, Val: 5},
	}, stack.Results[1])
	assert.Equal(t, pmark.Pair[int]{
		L: pmark.Cell[int]{Mark: pmark.BucketsCurlyL, Val: 3},
		R: pmark.Cell[int]{Mark: pmark.BucketsCurlyR, Val: 6},
	}, stack.Results[2])
}
