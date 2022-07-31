package tablety

type LineReader[TVal any] interface {

	// Read
	// returns io.EOF when finished
	// returns []TVal when empty
	Read() ([]TVal, error)
}
