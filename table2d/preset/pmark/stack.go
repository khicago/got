package pmark

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"

	"github.com/khicago/got/util/delegate"
)

type (
	// Cell - a cell with mark and payload
	Cell[TPayload constraints.Ordered] struct {
		Val  TPayload `json:"val"`
		Mark `json:"mark"`
	}

	// Pair - a pair of cells, generally used for pairing results
	Pair[TPayload constraints.Ordered] struct {
		L Cell[TPayload] `json:"l"`
		R Cell[TPayload] `json:"r"`
	}

	// Stack - a stack of cells, supporting mark pairings
	Stack[TPayload constraints.Ordered] struct {
		stack   []Cell[TPayload]
		Results []Pair[TPayload]
	}
)

var (
	ErrInsufficientStackLen = errors.New("insufficient stack length")
	ErrMarkAreNotRegistered = errors.New("mark are not registered")
)

// Between - check if the pair is between the given value
func (p Pair[TPayload]) Between(val TPayload, leftInclusive, rightInclusive bool) bool {
	if leftInclusive && rightInclusive {
		return p.L.Val <= val && val <= p.R.Val
	}
	if leftInclusive {
		return p.L.Val <= val && val < p.R.Val
	}
	if rightInclusive {
		return p.L.Val < val && val <= p.R.Val
	}
	return p.L.Val < val && val < p.R.Val
}

// Inside - check if the pair is inside the given value
func (p Pair[TPayload]) Inside(val TPayload) bool {
	return p.Between(val, false, false)
}

// NewStack - create a new stack
func NewStack[TPayload constraints.Ordered](cap int) *Stack[TPayload] {
	ret := Stack[TPayload]{}

	if cap > 0 {
		ret.stack = make([]Cell[TPayload], 0, cap)
		ret.Results = make([]Pair[TPayload], 0, (cap+1)>>1)
	} else {
		ret.stack = make([]Cell[TPayload], 0)
		ret.Results = make([]Pair[TPayload], 0o1)
	}
	return &ret
}

// Len returns the length of the stack
func (s *Stack[TPayload]) Len() int {
	return len(s.stack)
}

// Push
// only registered mark can be pushed into the stack
func (s *Stack[TPayload]) Push(mark Mark, val TPayload) error {
	if !mark.Registered() {
		return fmt.Errorf("%w, mark = %s", ErrMarkAreNotRegistered, mark)
	}
	s.stack = append(s.stack, Cell[TPayload]{Mark: mark, Val: val})
	return nil
}

// PeekMark returns the mark of top cell of the stack
func (s *Stack[TPayload]) PeekMark() (mark Mark, err error) {
	p, e := s.Peek()
	return p.Mark, e
}

// Peek returns the top cell of the stack
func (s *Stack[TPayload]) Peek() (cell Cell[TPayload], err error) {
	l := s.Len()
	if l == 0 {
		return Cell[TPayload]{Mark: NIL}, ErrInsufficientStackLen
	}
	return (s.stack)[l-1], nil
}

// Pop returns the top cell of the stack
func (s *Stack[TPayload]) Pop() (cell Cell[TPayload], err error) {
	l := s.Len()
	if l == 0 {
		return Cell[TPayload]{Mark: NIL}, ErrInsufficientStackLen
	}
	cell = (s.stack)[l-1]
	s.stack = (s.stack)[:l-1]
	return cell, nil
}

// Consume - try consumes the top cell of the stack
// if the top cell's mark is pair with the given mark, then the top cell will
// be popped out and paired with the given cell
// and the paired cell will be pushed into the result list
// otherwise, the given cell will be pushed into the stack
//
// actionPairing - action to be called when each input mark and payload are
// inputted, when the top cell's mark is pair with the given mark, the action
// will be called with true and the pair, and the given cell will be the right
// cell; otherwise, the action will be called with false and the pair, the given
// cell will be the left cell of the pair
func (s *Stack[TPayload]) Consume(mark Mark, payload TPayload,
	actionPairing delegate.Action2[bool, Pair[TPayload]],
) error {
	l := s.Len()
	if l == 0 {
		actionPairing.TryCall(false, Pair[TPayload]{L: Cell[TPayload]{Mark: mark, Val: payload}})
		return s.Push(mark, payload)
	}
	peak, err := s.PeekMark()
	if err != nil {
		return err
	}
	if IsPair(peak, mark) {
		lVal, err := s.Pop()
		if err != nil {
			return err
		}

		res := Pair[TPayload]{
			L: lVal,
			R: Cell[TPayload]{Mark: mark, Val: payload},
		}
		s.Results = append(s.Results, res)

		actionPairing.TryCall(true, res)
		return nil
	}
	actionPairing.TryCall(false, Pair[TPayload]{L: Cell[TPayload]{Mark: mark, Val: payload}})
	return s.Push(mark, payload)
}
