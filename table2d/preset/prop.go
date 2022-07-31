package preset

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/typer"
	"sort"
	"strconv"
	"strings"
)

type (
	Col = int

	PropTable map[int64]*Prop // PID => Prop

	// Prop
	// implementation of IProp
	Prop struct {
		p        map[Col]pseal.Seal
		children []MarkPair
	}

	MarkPair struct {
		BeginMarkCol Col
		EndMarkCol   Col
	}

	// List
	// implementation of IProp
	List []pseal.Seal

	IProp interface {
		Len() int
		Get(col Col) pseal.Seal
		Child(col Col) (IProp, error) // 直接把列表也放在 property 中, 方便按类型取
	}

	MarkPairProp struct {
		prop IProp
		MarkPair
		ValCols []Col
	}
)

var (
	_ IProp = &Prop{}
	_ IProp = &List{}
	_ IProp = &MarkPairProp{}
)

var (
	ErrSealFormatError = errors.New("seal marshal format error")
)

func NewProp() *Prop {
	return &Prop{
		p:        make(map[Col]pseal.Seal),
		children: make([]MarkPair, 0),
	}
}

func (p *Prop) MarshalJSON() ([]byte, error) {
	strs := make([]string, 0, len(p.p))
	keys := typer.Keys(p.p)
	sort.Ints(keys)
	for _, col := range keys {
		seal := (p.p)[col]
		strs = append(strs, fmt.Sprintf("%v:%v:%v", col, seal.Type().Name(), seal.Val()))
	}
	return json.Marshal(strs)
}

func (p *Prop) UnmarshalJSON(bytes []byte) error {
	strs := make([]string, 0)
	if err := json.Unmarshal(bytes, &strs); err != nil {
		return err
	}
	for _, str := range strs {
		values := strings.Split(str, ":")[0:]
		if len(values) < 3 {
			return ErrSealFormatError
		}
		col, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}
		ty := pseal.SymToType(values[1])
		s, err := ty.SealStr(str[len(values[0])+len(values[1])+2:])
		if err != nil {
			return err
		}
		p.p[Col(col)] = s
	}
	return nil
}

var (
	_ json.Marshaler   = &Prop{}
	_ json.Unmarshaler = &Prop{}
)

func propChild(prop IProp, col Col) (IProp, error) {
	seal := prop.Get(col)
	switch seal.Type() {
	case pseal.TyNil:
		return nil, ErrPropertyNil
	case pseal.TyAny:
		if v, ok := seal.Val().(IProp); ok {
			return v, nil
		}
		return nil, ErrPropertyType
	case pseal.TyMark:
		// todo
	}
	return nil, ErrPropertyType
}

func (p *Prop) Len() int {
	return len(p.p)
}

func (p *Prop) Get(col Col) pseal.Seal {
	v, ok := (p.p)[col]
	if !ok {
		return pseal.Nil
	}
	return v
}

func (p *Prop) Child(col Col) (IProp, error) {
	return propChild(p, col)
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
	return propChild(p, col)
}

func (m MarkPairProp) Len() int {
	return len(m.ValCols)
}

func (m MarkPairProp) Get(col Col) pseal.Seal {
	if !typer.SliceContains(m.ValCols, col) {
		return pseal.Nil
	}

	return m.prop.Get(col)
}

func (m MarkPairProp) Child(col Col) (IProp, error) {
	if !typer.SliceContains(m.ValCols, col) {
		return nil, ErrPropertyNil
	}

	return m.prop.Child(col)
}
