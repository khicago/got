package tablety

import (
	"github.com/khicago/got/util/delegate"
)

type (
	Col = int
	Row = int

	Table2DReader[TVal any] interface {
		MaxRow() int
		MaxCol() int

		/* query methods */

		// Get returns the val in given position
		Get(row Row, col Col) TVal

		// First used to find the first position in witch the val matches the given predicate
		// row by row, colum by colum
		First(pred delegate.Predicate[TVal]) (Row, Col)

		LineReader() LineReader[TVal]
	}
)
