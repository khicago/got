package preset

import (
	"encoding/json"
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
	"sort"
	"strconv"
	"strings"
)

type (

	// The ColHeader is a nested structure of ColMeta information tables
	//and is also optimised for queries by means of cached indexes etc. A
	//Preset has one and only one root ColHeader
	ColHeader struct {

		// Def handled ColMeta's information contained directly in
		// this ColHeader
		Def map[Col]*ColMeta `json:"def"`

		// Sub handle children of the meta table
		// why Col data needs to be structured:
		// - Although it is possible to structure the query in the col
		// header and prop, it does not solve the renaming problem.
		// In practice, especially in the case of lists with structures
		// inside them, there are problems with renaming, so it would
		// be too complicated to use query indexing to make it.
		Sub map[Col]*ColHeader `json:"sub"`

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

func NewColHeader() *ColHeader {
	return &ColHeader{
		Def:          make(map[Col]*ColMeta),
		nameColIndex: make(map[string]Col),
		Sub:          make(map[Col]*ColHeader),
	}
}

func (header *ColHeader) Set(col Col, colDef *ColMeta) *ColHeader {
	colDef.Col = col
	colDef.Constraint = strs.TrimLower(colDef.Constraint)
	header.Def[col] = colDef
	return header
}

func (header *ColHeader) ColOf(name string) Col {
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

func (header *ColHeader) Get(col Col) *ColMeta {
	return header.Def[col]
}

func (header *ColHeader) GetByName(name string) *ColMeta {
	col := header.ColOf(name)
	if col == InvalidCol {
		return nil
	}
	return header.Get(col)
}

var (
	_ json.Marshaler = &ColHeader{}
	//_ json.Unmarshaler = &ColHeader{}
)

// MarshalJSON
//
//	todo: 支持结构化
func (header *ColHeader) MarshalJSON() ([]byte, error) {
	type marshal struct {
		Metas []string
		Sub   map[Col]*ColHeader
	}

	strs := make([]string, 0, len(header.Def))
	keys := typer.Keys(header.Def)
	sort.Ints(keys)
	for _, col := range keys {
		meta := (header.Def)[col]
		strs = append(strs, meta.String())
	}
	return json.Marshal(marshal{
		Metas: strs,
		Sub:   header.Sub,
	})
}

func (header *ColHeader) UnmarshalJSON(bytes []byte) error {
	type marshal struct {
		Metas []string
		Sub   map[Col]*ColHeader
	}

	input := marshal{}
	if err := json.Unmarshal(bytes, &input); err != nil {
		return err
	}
	for _, str := range input.Metas {
		values := strings.Split(str, ":")[0:]
		if len(values) < 4 {
			return ErrSealFormatError
		}
		col, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}
		meta := &ColMeta{}
		err = meta.UnmarshalJSON([]byte(str))
		if err != nil {
			return err
		}

		header.Set(Col(col), meta)
	}
	return nil
}
