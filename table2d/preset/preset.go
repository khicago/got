package preset

import (
	"errors"
	"fmt"
	"strings"

	"github.com/khicago/got/table2d/preset/pcol"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/inlog"
)

type (
	Preset struct {
		Headline  *pcol.ColHeader
		PropTable *PropTable
	}
)

var (
	ErrPropertyNil  = errors.New("preset err: the property is nil")
	ErrPropertyType = errors.New("preset err: property type error")
)

func NewPreset() *Preset {
	return &Preset{
		Headline:  pcol.NewColHeader(),
		PropTable: &PropTable{},
	}
}

// QueryByCol returns the seal of the property by the column
func (p *Preset) QueryByCol(pid PresetID, col pcol.Col) pseal.Seal {
	seal, ok := (*p.PropTable)[pid]
	if !ok {
		return pseal.Invalid
	}
	return seal.Get(col)
}

// Query returns the seal of the property by the path
func (p *Preset) Query(pid PresetID, pth ...string) pseal.Seal {
	cm := p.Headline.GetByPth(pth...)
	inlog.Infof("got cm = %v\n", cm)
	if cm == nil {
		inlog.Debugf("got wrong path when parse %d, pth = %v", pid, pth)
		return pseal.Invalid
	}
	return p.QueryByCol(pid, cm.Col)
}

// QueryS is the same as Query, but the path is a string separated by "/"
func (p *Preset) QueryS(pid PresetID, pthStr string) pseal.Seal {
	return p.Query(pid, strings.Split(pthStr, "/")...)
}

// ForEachOfCol - for a given column, call the foreach function for each property
func (p *Preset) ForEachOfCol(col pcol.Col, fn delegate.Handler2[PresetID, pseal.Seal]) error {
	for pid, row := range *p.PropTable {
		if err := fn(pid, row.Get(col)); err != nil {
			return err
		}
	}
	return nil
}

// ForEach - for each property, try execute the foreach function
func (p *Preset) ForEach(fn delegate.Handler2[PresetID, IProp], orderly bool) error {
	return p.PropTable.ForEach(fn, orderly)
}

// Window - create a new preset by filtering the property table with given filters
// the filters is a map of column path and predicate
// the predicate is a function that takes a seal and returns a bool value
// when the predicate or path be nil, it will be ignored
//
// the head line of the new preset will be the same as the old one (pointer)
// the property table of the new preset will be a copy of the old one
func (p *Preset) Window(filters map[string]delegate.PredicateE[pseal.Seal] /* predicate */) (window *Preset, err error) {
	filtered := p.PropTable
	for f, fn := range filters {
		if fn == nil || f == "" {
			continue // ignore the nil predicate
		}

		if len(*filtered) == 0 {
			break // no need to filter
		}

		colMeta := p.Headline.GetByPth(f)
		if colMeta == nil {
			return nil, errors.New("got wrong path when filter")
		}
		col := colMeta.Col
		// filter the property table, a copy will be made by the filter method
		filtered, err = filtered.Filter(func(prop IProp) (bool, error) {
			seal := prop.Get(col)
			if seal == pseal.Invalid {
				return false, nil // ignore the invalid seal
			}
			ok, e := fn(seal)
			if e != nil {
				return false, fmt.Errorf("got error when filter: %v", e)
			}
			return ok, nil
		})
		if err != nil {
			return nil, err
		}
	}
	return &Preset{
		Headline:  p.Headline,
		PropTable: filtered,
	}, nil
}

// Retrieve - retrieve the property by the given id
//func (p *Preset) Retrieve(func(pid PresetID, indx *preset.Indx) error) (IProp, bool) {
//	panic("implement me")
//}
