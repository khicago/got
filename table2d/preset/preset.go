package preset

import (
	"github.com/khicago/got/table2d/preset/pseal"
)

type (
	Preset struct {
		Header *ColMetaTable
		*PropTable
	}
)

func NewPreset() *Preset {
	return &Preset{
		Header:    NewColHeader(),
		PropTable: &PropTable{},
	}
}

func (p *Preset) Query(pid int64, col Col) pseal.Seal {
	seal, ok := (*p.PropTable)[pid]
	if !ok {
		return pseal.Nil
	}
	return seal.Get(col)
}
