package preset

import (
	"context"
	"errors"
	"fmt"
	"github.com/khicago/got/internal/utils"
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

	headerRoot, marksStack, err2 := ParseHeader(preset.Headline, colPID, colLen-1,
		lineOfMeta, typer.SlicePadRight(lineColName, colLen, ""), typer.SlicePadRight(lineConstraint, colLen, ""))
	if err2 != nil {
		return nil, err2
	}

	// try load data values
	// 这个机制主要保证在没有 ColHeader 的时候, Props 自己能降级到支持平铺的结构访问
	childrenCols := make(PropChildIndex)
	for i, p := range marksStack.Results {
		childrenCols[p.LVal] = marksStack.Results[i]
	}

	inlog.Debugf("[READER] start parse data, got header %s", utils.MarshalIndentPrintAll(headerRoot))
	for line, err = read(); err == nil; line, err = read() {
		inlog.Debugf("read data line, %v, %v \n", line, typer.AssertNotNil(line))
		prop := NewProp()
		prop.childrenCols = childrenCols
		headerRoot.ForeachCol(func(colMeta *ColMeta) {
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

func ParseHeader(root *ColHeader, colFrom, colTo int, lineOfMeta []string, lineColName []string, lineConstraint []string) (
	*ColHeader, *pmark.Stack[Col], error) {

	inlog.Debugf("lineOfMeta(%d):\t\t %s\n", len(lineOfMeta), utils.MarshalPrintAll(lineOfMeta))
	inlog.Debugf("lineColName(%d):\t %s\n", len(lineColName), utils.MarshalPrintAll(lineColName))
	inlog.Debugf("lineConstraint(%d):\t %s\n", len(lineConstraint), utils.MarshalPrintAll(lineConstraint))

	headerStack := []*ColHeader{root}
	colPush := func(pairing bool, event pmark.Pair[Col]) {

		if !pairing {
			child := NewColHeader()
			typer.SliceLast(headerStack).Children[event.LVal] = ColHeaderChild{
				ColHeader: child,
				Pair:      event,
			}
			headerStack = append(headerStack, child)
			inlog.Debugf("------------ header stack in %#v, %v\n", event, headerStack)
			return
		}
		headerStack = headerStack[:len(headerStack)-1]
		inlog.Debugf("------------ header stack out %#v, %v\n", event, headerStack)
	}

	// try to generate header
	inlog.Debugf("colMax is %v\n", colTo)
	marksStack := pmark.NewStack[Col](colTo - colFrom)
	for c := colFrom; c <= colTo; c++ {
		sym := strs.TrimLower(lineOfMeta[c])
		ty := pseal.SymToType(sym)

		if ty == pseal.TyNil {
			inlog.Warnf("try parse col header sealTy of col %v skipped\n", c)
			continue
		}

		typer.SliceLast(headerStack).Set(c, NewColMeta(c, sym, lineColName[c], lineConstraint[c]))
		if ty == pseal.TyMark {
			err := marksStack.Consume(pmark.Mark(sym), c, colPush)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if len(headerStack) != 1 {
		return nil, nil, errors.New("parse header failed, stack are not cleared")
	}

	return headerStack[0], marksStack, nil
}
