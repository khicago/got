package table2d

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/khicago/got/annoparse"
	"github.com/khicago/got/table2d/tablety"
	"github.com/khicago/got/util/reflecto"
	"github.com/khicago/got/util/strs"
)

type (
	AnnoCol[TCol ~int] struct {
		Col    TCol
		Parser string
		Param  string
	}

	ParseOption struct {
		Tag   string
		Parse func(string) (any, error)
	}
)

func (AnnoCol[TCol]) TagName() string {
	return "tb2d"
}

var (
	once           = sync.Once{}
	AnnoTemplate   = &AnnoCol[int]{}
	AnnoTableCache annoparse.AnnoTable[*AnnoCol[int]]
)

// ParseByCol
// read all line from LineReader, and insert all items to outSlicePointer
func ParseByCol(outSlicePointer any, reader tablety.LineReader[string], option ...ParseOption) error {
	itr := func() (line any, err error) { return reader.Read() }

	elemType, err := reflecto.GetSliceElementType(outSlicePointer)
	if err != nil {
		return fmt.Errorf("invalid input, %w", err)
	}

	elemSpawner := reflecto.NewAnySpawnerFromType(elemType)
	mapper := func(in any) (any, error) {
		line, ok := in.([]string)
		if !ok {
			return nil, fmt.Errorf("invalid input %v", in)
		}
		v := elemSpawner.Spawn()

		if e := ParseLineByCol(v, line, option...); e != nil {
			return nil, e
		}
		return v, nil
	}
	csvReaderExitValidator := func(iv any, err error) (bool, error) {
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}

	return reflecto.Iterator(itr).WriteTo(outSlicePointer,
		reflecto.ItrMapper(mapper),
		reflecto.ItrExitValidator(csvReaderExitValidator),
	)
}

func ParseLineByCol(out any, line []string, option ...ParseOption) (err error) {
	if AnnoTableCache == nil {
		AnnoTableCache, err = annoparse.ExtractAnno(out, AnnoTemplate)
		if err != nil {
			return err
		}
	}

	if err = reflecto.ForEachField(out, func(fCtx reflecto.FieldContext) error {
		aCSV, ok := AnnoTableCache.Get(fCtx.Path)
		if !ok {
			return nil
		}

		valStr := line[aCSV.Col]
		var value any
		parser := aCSV.Parser
		switch {
		case "" == parser || strs.StartsWith(parser, "plain"):
			value, err = strs.Conv2PlainType(valStr, fCtx.Type)
		case strs.StartsWith(parser, "time"):
			if aCSV.Param == "" {
				value, err = time.Parse(time.RFC3339, valStr)
			} else {
				value, err = time.Parse(aCSV.Param, valStr)
			}
		default:
			for _, opt := range option {
				if aCSV.Parser == opt.Tag {
					value, err = opt.Parse(valStr)
					break
				}
			}
		}

		if err != nil {
			return err
		}
		fCtx.Value.Set(reflect.ValueOf(value))

		return nil
	},
		reflecto.ForEachFieldOptions.OnlyExported(),
		reflecto.ForEachFieldOptions.Drill(-1),
	); err != nil {
		return err
	}

	return nil
}
