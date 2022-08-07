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
	// ColMeta defines how this column is resolved, only the first Mark
	// of the substructure is recorded in the upper ColMeta list
	ColMeta struct {

		// Col is the column number of the meta, this information is also
		// available in the upper level table, the redundancy is provided
		// here mainly to pass the full information separately and later.
		Col

		// Type is the type of data for the column
		Type pseal.Type

		// Name column names are automatically converted to snake lower
		// case as the data is tiled, so even if there are substructures,
		// there can be no renamed columns
		Name string

		// Sym used to store the actual symbol filled in, to do some data
		// comparison, and to derive the type. When serialising and
		// de-serialising, the symbol should be serialised and de-serialised
		// in preference to the type.
		Sym string

		// Constraint is responsible for data validation, correlation,
		// processing, etc.
		Constraint string
	}

	// The ColMetaTable is a nested structure of ColMeta information tables
	//and is also optimised for queries by means of cached indexes etc. A
	//Preset has one and only one root ColMetaTable
	ColMetaTable struct {

		// Def handled ColMeta's information contained directly in
		// this ColMetaTable
		Def map[Col]*ColMeta

		// Sub handle children of the meta table
		// why Col data needs to be structured:
		// - Although it is possible to structure the query in the col
		// header and prop, it does not solve the renaming problem.
		// In practice, especially in the case of lists with structures
		// inside them, there are problems with renaming, so it would
		// be too complicated to use query indexing to make it.
		Sub map[Col]*ColMetaTable

		nameColIndex map[string]Col
	}

	//
	//IColMetaTable interface {
	//	Set(col Col, colDef *ColMeta) IColMetaTable
	//	ColOf(name string) Col
	//	Get(col Col) *ColMeta
	//	GetByName(name string) *ColMeta
	//}

)

func NewColHeader() *ColMetaTable {
	return &ColMetaTable{
		Def:          make(map[Col]*ColMeta),
		nameColIndex: make(map[string]Col),
		Sub:          make(map[Col]*ColMetaTable),
	}
}

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

// todo: 支持结构化
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

// todo: 支持结构化
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
