package preset

import (
	"context"
	"fmt"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/table2d/tablety"
	"github.com/khicago/got/util/strs"
)

type (
	// Raw
	// - 读入的原始行
	Raw      map[Col]string
	RawTable []Raw
)

// Read
// |@    |ID     |INT           |Float    |[         |       |      ]|[        |{|          |#        |          |#    |          |}|]|BOOL  |STRING          |
// |     |link(@)|test($>1,$<50)|test($%2)|link(item)|       |       |         | |link(wear)|         |link(wear)|     |link(wear)| | |      |enum(A|S|SR|SSR)|
// |     |LvUp   |Power         |Magic    |InitItems |       |       |InitStuff| |Weapon    |         |Shoes     |     |Hat       | | |IsHero|Gene            |
// |10001|10002  |12            |         |          |1010001|1010002|         | |1020001   |The Sword|1020101   |Speed|1030001   | | |Y     |S               |
func Read(ctx context.Context, reader tablety.Table2DReader[string]) *Preset {
	rowOfMeta, colPID := reader.First(func(val string) bool {
		return pseal.TyPID.SymMatch(val) != ""
	})
	if rowOfMeta < 0 {
		fmt.Printf("warn: cannot find pid\n")
		return nil
	}
	rowConstraint, rowColName, rowValueStart := rowOfMeta+1, rowOfMeta+2, rowOfMeta+3
	preset := NewPreset()

	for c := colPID; c <= reader.MaxCol(); c++ {
		sym := reader.Get(rowOfMeta, c)
		ty := pseal.SymToType(sym)
		if ty == pseal.TyNil {
			fmt.Printf("warn: try parse sealTy of row %v col %v failed\n", rowOfMeta, c)
			continue
		}
		constraint := reader.Get(rowConstraint, c)
		name := reader.Get(rowColName, c)

		preset.Header.Set(c, &ColMeta{
			Name:       strs.Conv2Snake(name),
			Type:       ty,
			Constraint: constraint,
		})
	}

	for r := rowValueStart; r <= reader.MaxRow(); r++ {
		prop, pid := Prop{}, int64(-1)
		for c, meta := range preset.Header.Def {
			str := reader.Get(r, c)
			val, err := meta.Type.SealByStr(str)
			if err != nil {
				fmt.Printf("warn: sealbystr of row %v col %v failed, ty= %v, str= %v, got err %v \n", r, c, meta.Type, str, err)
				continue
			}
			prop[c] = val
		}
		pid, err := prop.Get(colPID).PID()
		if err != nil {
			fmt.Printf("warn: read pid of row %v col %v failed, got err %v \n", r, colPID, err)
		}
		if pid == -1 {
			continue
		}
		(*preset.PropTable)[pid] = prop
	}

	return preset
}
