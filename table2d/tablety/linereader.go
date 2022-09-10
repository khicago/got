package tablety

import "io"

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

func (cl *commonLineReader[TVal]) Read() ([]TVal, error) {
	if cl.index >= len(cl.val) {
		return nil, io.EOF
	}
	s := cl.val[cl.index]
	cl.index++
	return s, nil
}

func Warp[TVal any](val [][]TVal) LineReader[TVal] {
	return &commonLineReader[TVal]{
		val: val,
	}
}
