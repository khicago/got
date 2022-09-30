package pmark

import (
	"errors"
	"fmt"
	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/typer"
)

type (
	Pair[TPayload any] struct {
		L, R       Mark
		LVal, RVal TPayload
	}

	Stack[TPayload any] struct {
		s       []Mark
		p       []TPayload
		Results []Pair[TPayload]
	}
)

var (
	ErrInsufficientStackLen = errors.New("insufficient stack length")
	ErrMarkAreNotRegistered = errors.New("mark are not registered")
)

func NewStack[TPayload any](cap int) *Stack[TPayload] {
	ret := Stack[TPayload]{}
	if cap > 0 {
		ret.s = make([]Mark, 0, cap)
		ret.p = make([]TPayload, 0, cap)
		ret.Results = make([]Pair[TPayload], 0, (cap+1)>>1)
	} else {
		ret.s = make([]Mark, 0)
		ret.p = make([]TPayload, 0)
		ret.Results = make([]Pair[TPayload], 0)
	}
	return &ret
}

func (s *Stack[TPayload]) Len() int {
	return len(s.s)
}

// Push
// only registered mark can be pushed into the stack
func (s *Stack[TPayload]) Push(mark Mark, payload TPayload) error {
	if !mark.Registered() {
		return fmt.Errorf("%s %w", mark, ErrMarkAreNotRegistered)
	}
	s.s = append(s.s, mark)
	s.p = append(s.p, payload)
	return nil
}

func (s *Stack[TPayload]) Peak() (mark Mark, err error) {
	l := s.Len()
	if l == 0 {
		return NIL, ErrInsufficientStackLen
	}
	return (s.s)[l-1], nil
}

func (s *Stack[TPayload]) Pop() (mark Mark, payload TPayload, err error) {
	l := s.Len()
	if l == 0 {
		return NIL, typer.ZeroVal[TPayload](), ErrInsufficientStackLen
	}
	mark, s.s = (s.s)[l-1], (s.s)[:l-1]
	payload, s.p = (s.p)[l-1], (s.p)[:l-1]
	return mark, payload, nil
}

func (s *Stack[TPayload]) Consume(mark Mark, payload TPayload, actionPairing delegate.Action2[bool, Pair[TPayload]]) error {
	l := s.Len()
	if l == 0 {
		actionPairing.TryCall(false, Pair[TPayload]{L: mark, LVal: payload})
		return s.Push(mark, payload)
	}
	peak, err := s.Peak()
	if err != nil {
		return err
	}
	if IsPair(peak, mark) {
		m0, p0, err := s.Pop()
		if err != nil {
			return err
		}

		res := Pair[TPayload]{
			L:    m0,
			R:    mark,
			LVal: p0,
			RVal: payload,
		}
		s.Results = append(s.Results, res)

		actionPairing.TryCall(true, res)
		return nil
	}
	actionPairing.TryCall(false, Pair[TPayload]{L: mark, LVal: payload})
	return s.Push(mark, payload)
}
