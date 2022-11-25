package pseal

import (
	"errors"
	"strconv"

	"github.com/khicago/got/util/strs"
)

var ErrPIDMustBeUnfilled = errors.New("pid must be filled")

func (ty Type) SealDefault() Seal { // todo: refine this
	return ty.Seal(ty.Default())
}

// Seal - seal a value into a Seal, with type check
// return Invalid if type check failed
func (ty Type) Seal(val any) Seal {
	if !ty.Assert(val) {
		return Invalid
	}
	return NewSeal(ty, val)
}

// SealStr - seal a string into a Seal, with type check
func (ty Type) SealStr(val string) (Seal, error) {
	val = strs.TrimLower(val)
	if val == "" {
		switch ty {

		case TyPID:
		case TyAny, TyID, TyBool, TyInt, TyFloat, TyMemo, TyMark: // 除了 PID 为空时都用默认 seal
			return ty.SealDefault(), nil
		}
		return Invalid, ErrPIDMustBeUnfilled
	}

	switch ty {
	case TyAny:
		return ty.Seal(val), nil
	case TyPID:
		pid, err := strconv.ParseInt(val, 10, 64)
		return ty.Seal(pid), err
	case TyID:
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ty.SealDefault(), nil
		}
		return ty.Seal(id), err
	case TyBool:
		return ty.Seal(val == "y" || val == "true"), nil
	case TyInt:
		i, err := strconv.Atoi(val)
		if err != nil {
			return ty.SealDefault(), nil
		}
		return ty.Seal(i), err
	case TyFloat:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return ty.SealDefault(), nil
		}
		return ty.Seal(f), err
	case TyString:
		return ty.Seal(val), nil
	case TyMemo:
		return ty.Seal(val), nil
	case TyMark:
		return ty.Seal(val), nil
	}
	return Invalid, nil
}

func Any(val any) Seal {
	return TyAny.Seal(val)
}

func PID(val int64) Seal {
	return TyPID.Seal(val)
}

func ID(val int64) Seal {
	return TyID.Seal(val)
}

func Bool(val bool) Seal {
	return TyBool.Seal(val)
}

func Int(val int) Seal {
	return TyInt.Seal(val)
}

func Float(val float64) Seal {
	return TyFloat.Seal(val)
}

func String(val string) Seal {
	return TyString.Seal(val)
}

func Memo(val string) Seal {
	return TyMemo.Seal(val)
}

func Mark(val string) Seal {
	return TyMark.Seal(val)
}
