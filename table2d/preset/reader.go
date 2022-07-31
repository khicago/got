package preset

import (
	"context"
	"errors"
	"fmt"
	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/table2d/tablety"
	"github.com/khicago/got/util/inlog"
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
	"io"
)

type (
	// Raw
	// - 读入的原始行
	Raw      map[Col]string
	RawTable []Raw
)

var (
	ErrPresetMarkError   = errors.New("preset mark error")
	ErrPresetFormatError = errors.New("preset format error")
)

func assertRead[TVal any](reader tablety.LineReader[TVal], validator typer.Predicate[[]TVal]) ([]TVal, error) {
	ln, err := reader.Read()
	inlog.Debugf("-- READ %#v %v\n", ln, err)
	if err != nil { // error occurred, maybe io.EOF
		return nil, err
	}
	if ln == nil || validator == nil { // ended
		return nil, nil
	}
	if !validator(ln) {
		return nil, ErrPresetFormatError
	}
	return ln, nil
}

// Read
// |@    |ID     |INT           |Float    |[         |       |      ]|[        |{|          |#        |          |#    |          |}|]|BOOL  |STRING          |
// |     |link(@)|test($>1,$<50)|test($%2)|link(item)|       |       |         | |link(wear)|         |link(wear)|     |link(wear)| | |      |enum(A|S|SR|SSR)|
// |     |LvUp   |Power         |Magic    |InitItems |       |       |InitStuff| |Weapon    |         |Shoes     |     |Hat       | | |IsHero|Gene            |
// |10001|10002  |12            |         |          |1010001|1010002|         | |1020001   |The Sword|1020101   |Speed|1030001   | | |Y     |S               |
func Read(ctx context.Context, reader tablety.Table2DReader[string]) (*Preset, error) {
	return ReadLines(ctx, reader.LineReader())
}

func ReadLines(ctx context.Context, reader tablety.LineReader[string]) (*Preset, error) {
	rowCount, colPID, colMax := 0, -1, 0

	// read returns nil, nil when finished
	read := func() (ln []string, err error) {
		for len(ln) == 0 {
			rowCount++ // start from 1
			ln, err = assertRead(reader, func(v []string) bool { return len(v) >= colMax })
			if err != nil {
				inlog.Debugf("- ln %d fin at, %#v, %s\n", rowCount, ln, err)
				if err == io.EOF {
					return nil, nil
				}
				return nil, fmt.Errorf("%w, row= %v", err, rowCount)
			}
		}
		inlog.Debug("- ln", rowCount, ln)
		return ln, nil
	}

	var lineOfMeta []string = nil
GetColPID:
	for line, err := read(); typer.AssertNotNil(line); line, err = read() {
		if err != nil {
			return nil, err
		}
		for c, v := range line {
			inlog.Debugf("sym pid got: %d, c %v sym %v\n", rowCount, c, v)
			if !typer.IsZero(pseal.TyPID.SymMatch(v)) {
				lineOfMeta = line
				colPID = c

				break GetColPID
			}
		}
	}

	if lineOfMeta == nil {
		return nil, ErrPresetMarkError
	}
	colMax = len(lineOfMeta) - 1

	lineConstraint, err := read()
	if err != nil {
		return nil, err
	}

	lineColName, err := read()
	if err != nil {
		return nil, err
	}

	marksStack := make([]Col, 0, colMax+1)
	preset := NewPreset()

	for c := colPID; c <= colMax; c++ {
		sym := strs.TrimLower(lineOfMeta[c])
		ty := pseal.SymToType(sym)

		if ty == pseal.TyNil {
			inlog.Warnf("warn: try parse col header sealTy of row %v col %v skipped\n", rowCount, c)
			continue
		}

		cMeta := &ColMeta{
			Type:       ty,
			Sym:        sym,
			Name:       strs.Conv2Snake(lineColName[c]),
			Constraint: lineConstraint[c],
		}
		preset.Header.Set(c, cMeta)

		if ty == pseal.TyMark {
			lStack := len(marksStack)
			if lStack == 0 {
				marksStack = append(marksStack, c)
				continue
			}
			peak := marksStack[lStack-1]
			if pmark.IsPair(preset.Header.Get(peak).Sym, sym) {
				marksStack = marksStack[:lStack-1]
			}
		}
	}

	l, e := read()
	for ; l != nil; l, e = read() {
		inlog.Debugf("line, %v \n", l)
		prop := NewProp()
		for c, meta := range preset.Header.Def {
			str := l[c]
			val, err := meta.Type.SealStr(str)
			if err != nil {
				inlog.Warnf("warn: sealbystr of row %v col %v failed, ty= %v, str= %v, got err %v \n", rowCount, c, meta.Type, str, err)
				continue
			}
			prop.p[c] = val
		}
		pid, err := prop.Get(colPID).PID()
		if err != nil {
			inlog.Warnf("warn: read pid of row %v col %v failed, got err %v \n", rowCount, colPID, err)
			continue
		}
		if pid == -1 {
			continue
		}
		(*preset.PropTable)[pid] = prop
	}
	if e != nil {
		return nil, e
	}

	return preset, nil
}
