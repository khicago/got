package tablety

import (
	"errors"
	"io"

	"github.com/khicago/got/util/delegate"
)

type (
	LineReader[TVal any] interface {
		// Read
		// returns io.EOF when finished
		// returns []TVal when empty
		Read() ([]TVal, error)
	}

	commonLineReader[TVal any] struct {
		index int
		val   [][]TVal
	}
)

var ErrValidateFailed = errors.New("validate failed")

func (cl *commonLineReader[TVal]) Read() ([]TVal, error) {
	if cl.index >= len(cl.val) {
		return nil, io.EOF
	}
	s := cl.val[cl.index]
	cl.index++
	return s, nil
}

func WarpLineReader[TVal any](val [][]TVal) LineReader[TVal] {
	return &commonLineReader[TVal]{
		val: val,
	}
}

func AssertRead[TVal any](reader LineReader[TVal], validator delegate.Predicate[[]TVal]) ([]TVal, error) {
	ln, err := reader.Read()
	if err != nil { // error occurred, maybe io.EOF
		return nil, err
	}
	if ln == nil || validator == nil { // finished
		return ln, nil
	}
	if !validator(ln) {
		return nil, ErrValidateFailed
	}
	return ln, nil
}
