package preset

import (
	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/typer"
)

type PropTable map[PresetID]*Prop // PID => Prop

func (pt *PropTable) Get(pid PresetID) IProp {
	return (*pt)[pid]
}

func (pt *PropTable) ForEach(fn delegate.Handler2[PresetID, IProp], orderly bool) error {
	if fn == nil {
		return nil
	}
	if orderly {
		keys := typer.KeysSorted(*pt)
		for _, pid := range keys {
			if err := fn(pid, pt.Get(pid)); err != nil {
				return err
			}
		}
		return nil
	}
	for pid, prop := range *pt {
		if err := fn(pid, prop); err != nil {
			return err
		}
	}
	return nil
}
