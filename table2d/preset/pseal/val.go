package pseal

import (
	"github.com/khicago/got/util/typer"
)

type (
	Seal struct {
		ty  Type
		val any
	}
)

var Nil = Seal{
	ty:  TyNil,
	val: DefaultAny,
}

func NewSeal(ty Type, val any) Seal {
	return Seal{
		ty:  ty,
		val: val,
	}
}

func (s Seal) Type() Type {
	return s.ty
}

func (s Seal) Val() any {
	return s.val
}

func (s Seal) PID() (int64, error) {
	return typer.Convert[int64](s.val, DefaultPID)
}

func (s Seal) ID() (int64, error) {
	i64, err := typer.Convert[int64](s.val, DefaultID)
	if err != nil {
		i, err2 := typer.Convert[int](s.val, int(DefaultID))
		if err2 != nil {
			return 0, err
		}
		i64 = int64(i)
	}
	return i64, nil
}

func (s Seal) Bool() (bool, error) {
	return typer.Convert[bool](s.val, DefaultBool)
}

func (s Seal) Int() (int, error) {
	return typer.Convert[int](s.val, DefaultInt)
}

func (s Seal) Float() (float64, error) {
	return typer.Convert[float64](s.val, DefaultFloat)
}

func (s Seal) String() (string, error) {
	return typer.Convert(s.val, DefaultString)
}

func (s Seal) Memo() (string, error) {
	return typer.Convert(s.val, DefaultMemo)
}
