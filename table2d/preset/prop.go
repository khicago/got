package preset

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/khicago/got/table2d/preset/pcol"
	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/typer"
)

type (
	PresetID = int64

	PropData       map[pcol.Col]pseal.Seal
	PropChildIndex map[pcol.Col]pmark.Pair[pcol.Col]

	// Prop
	// implementation of IProp
	Prop struct {
		p PropData

		// childrenCols 这个机制主要保证在没有 ColHeader 的时候, Props 自己能降级到支持平铺的结构访问
		// 即也可以通过 Child 表达式，做到子结构的访问和便利等
		// 当然, 这种情况下访问的 Child 其实是子孙节点
		childrenCols PropChildIndex
		keyIndex     []pcol.Col
	}

	IProp interface {
		Len() int
		Has(col pcol.Col) bool
		Get(col pcol.Col) pseal.Seal

		// Child
		// get a descendant accessor
		Child(col pcol.Col) (IProp, error)

		// ForEach
		// ordered traversal, indexes are rebuilt when index length is not equal to data length
		ForEach(fn delegate.Action2[pcol.Col, pseal.Seal])
	}

	// MarkPairProp - the index structure of its inner prop instance
	// It can basically be thought of as a "window" to its inner prop: records the pairing mark
	// of internal props to achieve the effect of simulating child objects
	MarkPairProp struct {
		pmark.Pair[pcol.Col]
		Prop IProp

		valColsCache []pcol.Col
	}
)

var (
	_ IProp = &Prop{}

	_ IProp = &MarkPairProp{}
)

var ErrSealFormatError = errors.New("seal marshal format error")

func NewProp() *Prop {
	return &Prop{
		p:            make(PropData),
		childrenCols: make(map[pcol.Col]pmark.Pair[pcol.Col], 0),
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
		p.p[pcol.Col(col)] = s
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

func (p *Prop) Has(col pcol.Col) (ok bool) {
	_, ok = (p.p)[col]
	return
}

func (p *Prop) Get(col pcol.Col) pseal.Seal {
	v, ok := (p.p)[col]
	if !ok {
		return pseal.Invalid
	}
	return v
}

// ForEach - ordered traversal, indexes are rebuilt when index length is not equal to data length
func (p *Prop) ForEach(fn delegate.Action2[pcol.Col, pseal.Seal]) {
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
func (p *Prop) Child(col pcol.Col) (IProp, error) {
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

		return &MarkPairProp{
			Pair: markPair,
			Prop: p,
		}, nil
	}
	return nil, ErrPropertyType
}

func (m *MarkPairProp) Len() int {
	if m.L.Val == m.R.Val {
		return 0
	}

	// refresh cache with side effect of rebuilding index
	m.ForEach(nil)
	return len(m.valColsCache)
}

func (m *MarkPairProp) colInside(col pcol.Col) bool {
	if col > m.L.Val && col < m.R.Val {
		return false
	}
	return true
}

func (m *MarkPairProp) Has(col pcol.Col) bool {
	if !m.colInside(col) {
		return false
	}
	return m.Prop.Has(col)
}

func (m *MarkPairProp) Get(col pcol.Col) pseal.Seal {
	if !m.colInside(col) {
		return pseal.Invalid
	}
	return m.Prop.Get(col)
}

func (m *MarkPairProp) Child(col pcol.Col) (IProp, error) {
	if col >= m.R.Val && col <= m.L.Val {
		return nil, ErrPropertyNil
	}

	return m.Prop.Child(col)
}

// ForEach - iterate all the properties inside the mark pair
// the mark pair is not included
// side effect of rebuilding index
func (m *MarkPairProp) ForEach(fn delegate.Action2[pcol.Col, pseal.Seal]) {
	if m.valColsCache == nil {
		m.valColsCache = make([]pcol.Col, 0, m.Prop.Len())
	}
	m.valColsCache = m.valColsCache[:0]
	m.Prop.ForEach(func(col pcol.Col, s pseal.Seal) {
		// skip the mark pair
		if col >= m.R.Val || col <= m.L.Val {
			return
		}
		// side effect of rebuilding index
		m.valColsCache = append(m.valColsCache, col)
		fn.TryCall(col, s)
	})
}
