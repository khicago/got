package preset

import (
	"context"
	"errors"
	"fmt"
	"github.com/khicago/got/internal/utils"
	"github.com/khicago/got/table2d/preset/pcol"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/table2d/tablety"
	"github.com/khicago/got/util/inlog"
	"github.com/khicago/got/util/typer"
	"io"
)

type (
	// Raw
	// - 读入的原始行
	Raw      map[pcol.Col]string
	RawTable []Raw
)

var (
	ErrPresetMarkError   = errors.New("preset mark error")
	ErrPresetFormatError = errors.New("preset format error")
)

// Read
// |@    |ID     |INT           |Float    |[         |       |      ]|[        |{|          |#        |          |#    |          |}|]|BOOL  |STRING          |
// |     |link(@)|test($>1,$<50)|test($%2)|link(item)|       |       |         | |link(wear)|         |link(wear)|     |link(wear)| | |      |enum(A|S|SR|SSR)|
// |     |LvUp   |Power         |Magic    |InitItems |       |       |InitStuff| |Weapon    |         |Shoes     |     |Hat       | | |IsHero|Gene            |
// |10001|10002  |12            |         |          |1010001|1010002|         | |1020001   |The Sword|1020101   |Speed|1030001   | | |Y     |S               |
func Read(ctx context.Context, reader tablety.Table2DReader[string]) (*Preset, error) {
	return ReadLines(ctx, reader.LineReader())
}

func ReadLines(ctx context.Context, reader tablety.LineReader[string]) (*Preset, error) {
	rowCount, colPID, colLen := 0, -1, 0

	// read returns nil, nil when finished
	read := func() (ln []string, err error) {
		rowCount++ // start from 1
		ln, err = tablety.AssertRead(reader, nil)
		// func(v []string) bool { return len(v) >= colLen }) // warning, not blocking
		if err != nil {
			inlog.Debugf("ln(%d) fin at, %#v, %s\n", rowCount, ln, err)
			if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("%w, row= %v", err, rowCount)
			}
			return nil, err
		}

		inlog.Debugf("ln(%d) line: %v ;\n", rowCount, ln)
		return ln, nil
	}

	var (
		lineOfMeta []string = nil
		line       []string
		err        error
	)

GetColPID:
	for line, err = read(); err == nil; line, err = read() {
		inlog.Debugf("- parse meta line, %v, %v\n", line, err)
		for c, v := range line {
			inlog.Debugf("! sym pid got: %d, c %v sym %v\n", rowCount, c, v)
			if !typer.IsZero(pseal.TyPID.SymMatch(v)) {
				lineOfMeta = line
				colPID = c
				break GetColPID
			}
		}
	}
	if err != nil && err != io.EOF {
		return nil, err
	}

	if lineOfMeta == nil {
		return nil, fmt.Errorf("%w, got empty meta row", ErrPresetMarkError)
	}
	colLen = len(lineOfMeta)

	lineConstraint, err := read()
	if err != nil {
		return nil, err
	}

	lineColName, err := read()
	if err != nil {
		return nil, err
	}

	preset := NewPreset()
	marksStack, err2 := preset.Headline.ParseHeader(colPID, colLen-1, lineOfMeta, lineColName, lineConstraint)
	if err2 != nil {
		return nil, err2
	}

	// try load data values
	// 这个机制主要保证在没有 ColHeader 的时候, Props 自己能降级到支持平铺的结构访问
	childrenCols := make(PropChildIndex)
	for i, p := range marksStack.Results {
		childrenCols[p.LVal] = marksStack.Results[i]
	}

	inlog.Debugf("[READER] start parse data, got header %s", utils.MarshalIndentPrintAll(preset.Headline))
	for line, err = read(); err == nil; line, err = read() {
		inlog.Debugf("read data line, %v, %v \n", line, typer.InNotNil(line))
		prop := NewProp()
		prop.childrenCols = childrenCols
		preset.Headline.ForeachCol(func(colMeta *pcol.ColMeta) {
			if colMeta.Col >= len(line) {
				inlog.Warnf("try parse col %d of prop row %v skipped, length %d is insufficient %d\n", colMeta, rowCount, len(line))
				return
			}
			str := line[colMeta.Col]
			val, err := colMeta.Type.SealStr(str)
			if err != nil {
				inlog.Warnf("seal_by_str of row %v col %v failed, str= %v, got err %v \n", rowCount, colMeta, str, err)
				return
			}
			prop.p[colMeta.Col] = val
		}, true)

		pid, err := prop.Get(colPID).PID()
		if err != nil {
			inlog.Warnf("read pid of row %v col %v failed, got err %v \n", rowCount, colPID, err)
			continue
		}
		if pid == -1 {
			continue
		}
		(*preset.PropTable)[pid] = prop
	}
	if err != nil && err != io.EOF {
		return nil, err
	}

	return preset, nil
}
