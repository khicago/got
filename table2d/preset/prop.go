package preset

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/typer"
	"sort"
	"strconv"
	"strings"
)

type (
	Col = int

	PropTable map[int64]*Prop // PID => Prop

	PropData       map[Col]pseal.Seal
	PropChildIndex map[Col]pmark.Pair[Col]

	// Prop
	// implementation of IProp
	Prop struct {
		p PropData

		// childrenCols 这个机制主要保证在没有 ColHeader 的时候, Props 自己能降级到支持平铺的结构访问
		// 即也可以通过 Child 表达式，做到子结构的访问和便利等
		// 当然, 这种情况下访问的 Child 其实是子孙节点
		childrenCols PropChildIndex
		keyIndex     []Col
	}

	IProp interface {
		Len() int
		Has(col Col) bool
		Get(col Col) pseal.Seal

		// Child
		// get a descendant accessor
		Child(col Col) (IProp, error)

		// ForEach
		// ordered traversal, indexes are rebuilt when index length is not equal to data length
		ForEach(fn typer.DelegateAction2[Col, pseal.Seal])
	}

	MarkPairProp struct {
		pmark.Pair[Col]
		prop    IProp
		ValCols []Col
	}
)

var (
	_ IProp = &Prop{}

	_ IProp = &MarkPairProp{}
)

var (
	ErrSealFormatError = errors.New("seal marshal format error")
)

func NewProp() *Prop {
	return &Prop{
		p:            make(PropData),
		childrenCols: make(map[Col]pmark.Pair[Col], 0),
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

func (p *Prop) Len() int {
	return len(p.p)
}

func (p *Prop) Has(col Col) (ok bool) {
	_, ok = (p.p)[col]
	return
}

func (p *Prop) Get(col Col) pseal.Seal {
	v, ok := (p.p)[col]
	if !ok {
		return pseal.Nil
	}
	return v
}

func (p *Prop) ForEach(fn typer.DelegateAction2[Col, pseal.Seal]) {
	if len(p.keyIndex) != p.Len() {
		keys := typer.Keys(p.p)
		sort.Ints(keys)
		p.keyIndex = keys
	}

	for _, col := range p.keyIndex {
		fn.TryCall(col, p.p[col])
	}
}

// Child
// In fact, the current implementation can be described as a `descendant`
func (p *Prop) Child(col Col) (IProp, error) {
	seal := p.Get(col)
	switch seal.Type() {
	case pseal.TyNil:
		return nil, ErrPropertyNil
	case pseal.TyAny:
		if v, ok := seal.Val().(IProp); ok {
			return v, nil
		}
		return nil, ErrPropertyType
	case pseal.TyMark:
		markPair, ok := p.childrenCols[col]
		if !ok {
			if p.Has(col) {
				return nil, ErrPropertyType
			}
			return nil, ErrPropertyNil
		}

		return MarkPairProp{
			Pair: markPair,
			prop: p,
		}, nil
	}
	return nil, ErrPropertyType
}

func (m MarkPairProp) Len() int {
	return len(m.ValCols)
}

func (m MarkPairProp) colInside(col Col) bool {
	if col > m.LVal && col < m.RVal {
		return false
	}
	return true
}

func (m MarkPairProp) Has(col Col) bool {
	if !m.colInside(col) {
		return false
	}
	return m.prop.Has(col)
}

func (m MarkPairProp) Get(col Col) pseal.Seal {
	if !m.colInside(col) {
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

func (m MarkPairProp) ForEach(fn typer.DelegateAction2[Col, pseal.Seal]) {
	m.prop.ForEach(func(c Col, s pseal.Seal) {
		if c <= m.LVal || c >= m.RVal {
			return
		}
		fn.TryCall(c, s)
	})
}
