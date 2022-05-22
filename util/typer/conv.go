package typer

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrPropertyType = errors.New("got/typer: convert type failed")
)

func Convert[T any](val any, defaultV T) (T, error) {
	vv, ok := val.(T)
	if !ok {
		return defaultV, ErrPropertyType
	}
	return vv, nil
}

func ConvI2I64Any(val any) (int64, error) {
	switch t := val.(type) {
	case int:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case uint:
		return int64(t), nil
	case uint8:
		return int64(t), nil
	case uint16:
		return int64(t), nil
	case uint32:
		return int64(t), nil
	case uint64:
		return int64(t), nil
	case uintptr:
		return int64(t), nil
	case int64:
		return t, nil
	}
	return strconv.ParseInt(fmt.Sprintf("%v", val), 10, 16)
}
