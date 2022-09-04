package preset

import (
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/inlog"
	"strings"
)

type (
	Preset struct {
		Headline  *ColHeader
		PropTable *PropTable
	}
)

func NewPreset() *Preset {
	return &Preset{
		Headline:  NewColHeader(),
		PropTable: &PropTable{},
	}
}

func (p *Preset) QueryByCol(pid int64, col Col) pseal.Seal {
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
