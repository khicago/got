package typer

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

var ErrPropertyType = errors.New("got/typer: convert type failed")

func Convert[T any](val any, defaultV T) (T, error) {
	vv, ok := val.(T)
	if !ok {
		return defaultV, fmt.Errorf("%w, cannot convert %v(type %T) to type %T", ErrPropertyType, val, val, defaultV)
	}
	return vv, nil
}

func ConvertMust[T any](val any) T {
	zero := ZeroVal[T]()
	v, err := Convert[T](val, zero)
	if err != nil {
		panic(fmt.Errorf("conver to type %T failed, %w", zero, err))
	}
	return v
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

func I2Str[T constraints.Integer](num T) string {
	return strconv.FormatInt(int64(num), 10)
}

func Any[T any](v T) any {
	return v
}

func S2IMust[T constraints.Integer](str string) T {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(fmt.Errorf("s2i(%s) %T failed, %w", str, ZeroVal[T](), err))
	}
	return T(i)
}
