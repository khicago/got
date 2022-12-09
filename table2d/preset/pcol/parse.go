package pcol

import (
	"errors"

	"github.com/khicago/got/internal/utils"
	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/inlog"
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
)

func (header *ColHeader) ParseHeader(colFrom, colTo int, metaSymbols, names, constraints []string) (*pmark.Stack[Col], error) {
	inlog.Debugf("lineOfMeta(%d):\t\t %s\n", len(metaSymbols), utils.MarshalPrintAll(metaSymbols))
	inlog.Debugf("lineColName(%d):\t %s\n", len(names), utils.MarshalPrintAll(names))
	inlog.Debugf("lineConstraint(%d):\t %s\n", len(constraints), utils.MarshalPrintAll(constraints))

	headerStack := []*ColHeader{header}
	colPush := func(pairing bool, event pmark.Pair[Col]) {
		peak := typer.SliceLast(headerStack)
		// if not pairing, create a new child, and push it into stack
		if !pairing {
			child := &ColHeader{
				ColHeaderData: peak.ColHeaderData,
				Pair:          &event,
			}
			headerStack = append(headerStack, child)
			peak.Children[event.L.Val] = child
			inlog.Debugf("------------ header stack in %#v, %v\n", event, headerStack)
			return
		}
		// if pairing, overwrite the current header's pair
		peak.Pair = &event
		// pop the stack
		headerStack = headerStack[:len(headerStack)-1]
		inlog.Debugf("------------ header stack out %#v, %v\n", event, headerStack)
	}

	// try to generate header
	inlog.Debugf("colMax is %v\n", colTo)
	marksStack := pmark.NewStack[Col](colTo - colFrom)
	for c := colFrom; c <= colTo; c++ {
		sym := strs.TrimLower(metaSymbols[c])
		ty := pseal.SymToType(sym)

		if ty == pseal.TyNil {
			inlog.Warnf("try parse col header sealTy of col %v skipped\n", c)
			continue
		}

		typer.SliceLast(headerStack).Set(c, NewColMeta(c, sym,
			typer.SliceTryGet(names, c, ""),
			typer.SliceTryGet(constraints, c, ""),
		))

		// only mark will be pushed into stack
		if ty == pseal.TyMark {
			err := marksStack.Consume(pmark.Mark(sym), c, colPush)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(headerStack) != 1 {
		return nil, errors.New("parse header failed, stack are not cleared")
	}

	return marksStack, nil
}
