package preset

import (
	"github.com/khicago/got/table2d/preset/pseal"
)

type (
	Col = int

	PropTable map[int64]Prop // PID => Prop

	// Prop
	// implementation of IProp
	Prop map[Col]pseal.Seal

	// List
	// implementation of IProp
	List []pseal.Seal

	IProp interface {
		Len() int
		Get(col Col) pseal.Seal
		Child(col Col) (IProp, error) // 直接把列表也放在 property 中, 方便按类型取
	}
)

var _ IProp = &Prop{}
var _ IProp = &List{}

func (p *Prop) Len() int {
	return len(*p)
}

func (p *Prop) Get(col Col) pseal.Seal {
	v, ok := (*p)[col]
	if !ok {
		return pseal.Nil
	}
	return v
}

func (p *Prop) Child(col Col) (IProp, error) {
	seal := p.Get(col)
	if seal.Type() == pseal.TyNil {
		return nil, ErrPropertyNil
	}
	obj := seal.Val()
	if v, ok := obj.(IProp); ok {
		return v, nil
	}
	return nil, ErrPropertyType
}

func (p *List) Len() int {
	return len(*p)
}

func (p *List) Get(col Col) pseal.Seal {
	if col < 0 || col > p.Len() {
		return pseal.Nil
	}
	return (*p)[col]
}

func (p *List) Child(col Col) (IProp, error) {
	seal := p.Get(col)
	if seal.Type() == pseal.TyNil {
		return nil, ErrPropertyNil
	}
	obj := seal.Val()
	if v, ok := obj.(IProp); ok {
		return v, nil
	}
	return nil, ErrPropertyType
}
