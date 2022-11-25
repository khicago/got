package preset

import (
	"errors"
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
func (p *Preset) ForEachOfCol(col pcol.Col, foreach delegate.Handler2[PresetID, pseal.Seal]) error {
	for pid, row := range *p.PropTable {
		if err := foreach(pid, row.Get(col)); err != nil {
			return err
		}
	}
	return nil
}
