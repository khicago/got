package preset

import (
	"github.com/khicago/got/table2d/preset/pseal"
)

type (
	Preset struct {
		H *ColHeader
		*PropTable
	}
)

func NewPreset() *Preset {
	return &Preset{
		H:         NewColHeader(),
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

func (p *Preset) QueryS(pid int64, path []Col) pseal.Seal {
	// todo
	panic("to implement this")
}
