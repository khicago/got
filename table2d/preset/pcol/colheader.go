package pcol

import (
	"encoding/json"
	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/util/inlog"
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
	"sort"
	"strconv"
	"strings"
)

type (
	ColHeaderChild struct {
		*ColHeader
		pmark.Pair[Col]
	}

	ColHeaderChildren map[Col]ColHeaderChild

	// The ColHeader is a nested structure of ColMeta information tables
	//and is also optimised for queries by means of cached indexes etc. A
	//Preset has one and only one root ColHeader
	ColHeader struct {

		// Def handled ColMeta's information contained directly in
		// this ColHeader
		Def map[Col]*ColMeta `json:"def"`

		// Children handle children of the meta table
		// why Col data needs to be structured:
		// - Although it is possible to structure the query in the col
		// header and prop, it does not solve the renaming problem.
		// In practice, especially in the case of lists with structures
		// inside them, there are problems with renaming, so it would
		// be too complicated to use query indexing to make it.
		Children ColHeaderChildren `json:"sub"`

		nameColIndex map[string]Col
	}

	//
	//IColMetaTable interface {
	//	Set(col Col, colDef *ColMeta) IColMetaTable
	//	ColOf(name string) Col
	//	Get(col Col) *ColMeta
	//	GetByPth(name string) *ColMeta
	//}

)

var (
	_ json.Marshaler   = &ColHeader{}
	_ json.Unmarshaler = &ColHeader{}
)

func NewColHeader() *ColHeader {
	return &ColHeader{
		Def:      make(map[Col]*ColMeta),
		Children: make(ColHeaderChildren),

		nameColIndex: make(map[string]Col),
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

func (header *ColHeader) GetByIndex(index int) *ColMeta {
	keys := typer.Keys(header.Def)
	sort.Ints(keys)

	return header.Def[keys[index]]
}

func (hChild *ColHeaderChild) GetByMark(indexOrCol string) *ColMeta {
	if hChild.Pair.L == "[" {
		ind, err := strconv.Atoi(indexOrCol)
		if err != nil {
			return nil
		}
		return hChild.GetByIndex(ind)
	}
	col := hChild.ColOf(indexOrCol) // be mark for child
	if col == InvalidCol {
		return nil
	}
	return hChild.Get(col)
}

func (header *ColHeader) GetByPth(pth ...string) *ColMeta {
	var (
		recursive = ColHeaderChild{
			Pair: pmark.Pair[Col]{
				L: "{",
				R: "}",
			},
			ColHeader: header,
		}
		node *ColMeta
	)
	for i := range pth {
		node = recursive.GetByMark(pth[i])
		inlog.Debugf("-GetByPth> recursive node (i=%d),\tpth=%v,\tnode=`%v`\n", i, pth, node)
		if node == nil || i == len(pth)-1 {
			inlog.Debugf("-GetByPth> recursive node (final),\tpth=%v,\tnode=`%v`\n", pth, node)
			break
		}

		child, ok := recursive.Children[node.Col]
		if !ok {
			return nil
		}
		recursive = child

	}
	return node
}

func (header *ColHeader) ForeachCol(action func(colMeta *ColMeta), includeChildren bool) {
	for _, def := range header.Def {
		action(def)
	}
	if !includeChildren {
		return
	}
	for _, child := range header.Children {
		child.ForeachCol(action, true)
	}
}

// MarshalJSON
//
//	todo: 支持结构化
func (header *ColHeader) MarshalJSON() ([]byte, error) {
	type marshal struct {
		Metas []string
		Sub   ColHeaderChildren
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
		Sub:   header.Children,
	})
}

func (header *ColHeader) UnmarshalJSON(bytes []byte) error {
	type marshal struct {
		Metas []string
		Sub   ColHeaderChildren
	}

	input := marshal{}
	if err := json.Unmarshal(bytes, &input); err != nil {
		return err
	}
	for _, str := range input.Metas {
		values := strings.Split(str, ":")[0:]
		if len(values) < 4 {
			return ErrColHeaderUnMarshalFmtFailed
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
