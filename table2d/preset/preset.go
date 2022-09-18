package preset

import (
	"errors"
	"github.com/khicago/got/table2d/preset/pcol"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/inlog"
	"strings"
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

func (p *Preset) QueryByCol(pid int64, col pcol.Col) pseal.Seal {
	seal, ok := (*p.PropTable)[pid]
	if !ok {
		return pseal.Nil
	}
	return seal.Get(col)
}

func (p *Preset) Query(pid int64, pth ...string) pseal.Seal {
	cm := p.Headline.GetByPth(pth...)
	if cm == nil {
		inlog.Debugf("got wrong path when parse %d, pth = %v", pid, pth)
		return pseal.Nil
	}
	return p.QueryByCol(pid, cm.Col)
}

func (p *Preset) QueryS(pid int64, pthStr string) pseal.Seal {
	return p.Query(pid, strings.Split(pthStr, "/")...)
}
