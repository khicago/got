package typer

import "errors"

var (
	ErrPropertyType = errors.New("got/typer: convert type  failed")
)

func Convert[T any](val any, defaultV T) (T, error) {
	vv, ok := val.(T)
	if !ok {
		return defaultV, ErrPropertyType
	}
	return vv, nil
}
