package pseal

import (
	"fmt"
	"github.com/khicago/got/util/strs"
	"strconv"
)

func (ty Type) SealDefault() Seal { // todo: refine this
	return ty.Seal(ty.Default())
}

// Seal
// 返回 Nil 表示封装失败
func (ty Type) Seal(val any) Seal {
	if !ty.Assert(val) {
		return Nil
	}
	return NewSeal(ty, val)
}

func (ty Type) SealByStr(val string) (Seal, error) {
	val = strs.TrimLower(val)
	switch ty {
	case TyAny:
		return ty.Seal(val), nil
	case TyPID:
		fmt.Println("seal pid", val)
		pid, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ty.SealDefault(), err
		}
		return ty.Seal(pid), err
	case TyID:
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ty.SealDefault(), err
		}
		return ty.Seal(id), err
	case TyBool:
		return ty.Seal(val == "y" || val == "true"), nil
	case TyInt:
		i, err := strconv.Atoi(val)
		if err != nil {
			return ty.SealDefault(), err
		}
		return ty.Seal(i), err
	case TyFloat:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return ty.SealDefault(), err
		}
		return ty.Seal(f), err
	case TyString:
		return ty.Seal(val), nil
	case TyMemo:
		return ty.Seal(val), nil
	case TyMark:
		return ty.Seal(val), nil
	}
	return Nil, nil
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
