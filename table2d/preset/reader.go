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
		rowCount++ // start from 1
		ln, err = assertRead(reader, func(v []string) bool { return len(v) >= colMax })
		if err != nil {
			inlog.Debugf("- ln %d fin at, %#v, %s\n", rowCount, ln, err)
			if err == io.EOF {
				return nil, nil
			}
			return nil, fmt.Errorf("%w, row= %v", err, rowCount)
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

	preset := NewPreset()

	headerRoot, marksStack, err2 := ParseHeader(preset.H, colPID, colMax, lineOfMeta, lineColName, lineConstraint)
	if err2 != nil {
		return nil, err2
	}
	// try load data values

	childrenCols := make(PropChildIndex)
	for i, p := range marksStack.Results {
		childrenCols[p.LVal] = marksStack.Results[i]
	}

	inlog.Debugf("- header root %#v", headerRoot)
	l, e := read()
	for ; l != nil; l, e = read() {
		inlog.Debugf("line, %v \n", l)
		prop := NewProp()
		prop.childrenCols = childrenCols
		for c, meta := range headerRoot.Def {
			str := l[c]
			val, err := meta.Type.SealStr(str)
			if err != nil {
				inlog.Warnf("seal_by_str of row %v col %v failed, ty= %v, str= %v, got err %v \n", rowCount, c, meta.Type, str, err)
				continue
			}
			prop.p[c] = val
		}
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
	if e != nil {
		return nil, e
	}

	return preset, nil
}

func ParseHeader(root *ColHeader, colFrom, colTo int, lineOfMeta []string, lineColName []string, lineConstraint []string) (
	*ColHeader, *pmark.Stack[Col], error) {

	headerStack := []*ColHeader{root}
	colPush := func(pairing bool, event pmark.Pair[Col]) {
		if !pairing {
			newHeader := NewColHeader()
			typer.SliceLast(headerStack).Sub[event.LVal] = newHeader
			headerStack = append(headerStack, newHeader)

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

		typer.SliceLast(headerStack).Set(c, NewColMeta(c, sym, lineColName[c], lineConstraint[c])) // todo: 第一个 Mark 留在父结构里
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
