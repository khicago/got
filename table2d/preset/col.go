package preset

import (
	"encoding/json"
	"fmt"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
	"sort"
	"strconv"
	"strings"
)

type (
	ColMeta struct {
		Col
		Type       pseal.Type
		Name       string
		Sym        string
		Constraint string
	}
	//
	//IColMetaTable interface {
	//	Set(col Col, colDef *ColMeta) IColMetaTable
	//	ColOf(name string) Col
	//	Get(col Col) *ColMeta
	//	GetByName(name string) *ColMeta
	//}

	ColMetaTable struct {
		Def          map[Col]*ColMeta
		nameColIndex map[string]Col
	}
)

func NewColHeader() *ColMetaTable {
	return &ColMetaTable{
		Def:          make(map[Col]*ColMeta),
		nameColIndex: make(map[string]Col),
	}
}

//
//var _ IColMetaTable = &ColMetaTable{}

func (header *ColMetaTable) Set(col Col, colDef *ColMeta) *ColMetaTable {
	colDef.Col = col
	colDef.Constraint = strs.TrimLower(colDef.Constraint)
	header.Def[col] = colDef
	return header
}

func (header *ColMetaTable) ColOf(name string) Col {
	if v, ok := header.nameColIndex[name]; ok {
		return v
	}
	for k, v := range header.Def {
		if v.Name == name {
			header.nameColIndex[name] = k
			return k
		}
	}
	return InvalidCol
}

func (header *ColMetaTable) Get(col Col) *ColMeta {
	return header.Def[col]
}

func (header *ColMetaTable) GetByName(name string) *ColMeta {
	col := header.ColOf(name)
	if col == InvalidCol {
		return nil
	}
	return header.Get(col)
}

var (
	_ json.Marshaler   = &ColMetaTable{}
	_ json.Unmarshaler = &ColMetaTable{}
)

func (header *ColMetaTable) MarshalJSON() ([]byte, error) {
	strs := make([]string, 0, len(header.Def))
	keys := typer.Keys(header.Def)
	sort.Ints(keys)
	for _, col := range keys {
		meta := (header.Def)[col]
		strs = append(strs, fmt.Sprintf("%v:%s:%v:%s", col, meta.Sym, meta.Name, meta.Constraint))
	}
	return json.Marshal(strs)
}

func (header *ColMetaTable) UnmarshalJSON(bytes []byte) error {
	strs := make([]string, 0)
	if err := json.Unmarshal(bytes, &strs); err != nil {
		return err
	}
	for _, str := range strs {
		values := strings.Split(str, ":")[0:]
		if len(values) < 4 {
			return ErrSealFormatError
		}
		col, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}

		header.Set(Col(col), &ColMeta{
			Name:       values[2],
			Sym:        values[1],
			Type:       pseal.SymToType(values[1]),
			Constraint: str[len(values[0])+len(values[1])+len(values[2])+3:],
		})
	}
	return nil
}
