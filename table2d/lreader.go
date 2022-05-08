package table2d

type LineReader interface {
	Read() ([]string, error)
}
