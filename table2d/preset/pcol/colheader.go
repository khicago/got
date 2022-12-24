package pcol

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/util/inlog"
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
)

type (
	ColHeader struct {
		*ColHeaderData

		// ColHeaderChildren is the range of ColHeader
		*pmark.Pair[Col]

		// Children handle children of the meta table
		// why Col data needs to be structured:
		// - Although it is possible to structure the query in the col
		// header and prop, it does not solve the renaming problem.
		// In practice, especially in the case of lists with structures
		// inside them, there are problems with renaming, so it would
		// be too complicated to use query indexing to make it.
		Children ColHeaderChildren `json:"sub"`

		// nameColIndex is a cache of the name of the ColMeta
		nameColIndex map[string]Col
	}

	ColHeaderChildren map[Col]*ColHeader

	// The ColHeaderData is a nested structure of ColMeta information tables
	// and is also optimised for queries by means of cached indexes etc. A
	// Preset has one and only one root ColHeader
	ColHeaderData struct {
		// Def handled ColMeta's information contained directly in
		// this ColHeader
		Def map[Col]*ColMeta `json:"def"`
	}
)

var (
	_ json.Marshaler   = NewColHeader()
	_ json.Unmarshaler = NewColHeader()
)

// NewColHeader creates a new ColHeader
func NewColHeader() *ColHeader {
	return &ColHeader{
		ColHeaderData: &ColHeaderData{
			Def: make(map[Col]*ColMeta),
		},
		Children:     make(ColHeaderChildren),
		nameColIndex: make(map[string]Col),
	}
}

// ForkChild creates a new ColHeader with the same ColHeaderData
// a new pair should be set to the Children
// if you fork a new header of the child with a nil pair, then
// you got a root header
func (header *ColHeader) ForkChild(pair *pmark.Pair[Col]) *ColHeader {
	// for a child, the pair is a sub-range of the parent
	return &ColHeader{
		ColHeaderData: header.ColHeaderData,
		nameColIndex:  make(map[string]Col),
		Children:      make(ColHeaderChildren),
		Pair:          pair,
	}
}

// GetChildByCol returns the ColHeader of the child by the left col, the
// Children info is stored in the ColHeader.Children in the ParseHeader
// process
// returns nil if not found
func (header *ColHeader) GetChildByCol(childLeftCol Col) *ColHeader {
	return header.Children[childLeftCol]
}

// Set sets a ColMeta to the ColHeader
// copy the ColMeta to avoid the pointer being modified and
func (header *ColHeader) Set(col Col, colDef *ColMeta) *ColHeader {
	c := *colDef
	c.Col = col
	c.Constraint = strs.TrimLower(colDef.Constraint)
	header.Def[col] = &c
	return header
}

// ColOf returns the Col of the ColMeta by name
// returns InvalidCol if not found
func (header *ColHeader) ColOf(name string) Col {
	if v, ok := header.nameColIndex[name]; ok {
		return v
	}

	for col, v := range header.Def {
		if !header.IsSelfFiled(col, true) {
			continue
		}
		if v.Name == name {
			if exist, ok := header.nameColIndex[v.Name]; ok && exist != col {
				// todo: name conflict issue
				inlog.Warnf("name conflict: %s, %d, %d", v.Name, exist, col)
			}
			return col
		}
		header.nameColIndex[v.Name] = col
	}
	return InvalidCol
}

// Get returns the ColMeta by Col
// returns nil if not found
func (header *ColHeader) Get(col Col) *ColMeta {
	return header.Def[col]
}

// GetByIndex returns the ColMeta by index
// returns nil if not found
// e.g.: for
func (header *ColHeader) GetByIndex(index int) *ColMeta {
	var keys []Col
	for col := range header.Def {
		// child col is counted as only one col (as its left mark)
		if !header.IsSelfFiled(col, true) {
			continue
		}
		keys = append(keys, col)
	}
	sort.Ints(keys)

	return header.Def[keys[index]]
}

// GetByIndexOrColName returns the ColMeta by index or name
// if the ColHeader is a list, it will be treated as an index
// otherwise it will be treated as a name
func (header *ColHeader) GetByIndexOrColName(indexOrCol string) *ColMeta {
	// if the ColHeader is a list, it will be treated as an index
	// if the Pair is nil, it means that the ColHeader is a map
	if header.Pair != nil && header.Pair.L.Mark == "[" {
		ind, err := strconv.Atoi(indexOrCol)
		if err != nil {
			return nil
		}

		return header.GetByIndex(ind)
	}
	colName := indexOrCol
	col := header.ColOf(colName) // be name for child
	if col == InvalidCol {
		return nil
	}
	return header.Get(col)
}

// LenOf returns the length of a col
// if the col is not found, it will return -2
// if the col is a list, it will return the length of the list
// otherwise, it will return -1
func (header *ColHeader) LenOf(col Col) int {
	node := header.Get(col)
	if node == nil {
		return -2
	}
	if node.Sym == "[" {
		// get the children by node.Col
		count, child := 0, header.Children[node.Col]
		// count (child structure are counted as one col)
		for col := range child.Def {
			// child col is counted as only one col (as its left mark)
			if !child.IsSelfFiled(col, true) {
				continue
			}
			count++
		}
		return count
	}
	return -1
}

func (header *ColHeader) GetByPth(pth ...string) *ColMeta {
	var node *ColMeta
	recursive := header
	for i := range pth {
		node = recursive.GetByIndexOrColName(pth[i])
		inlog.Debugf("- GetByPth > recursive(i= %d),\tpth=%v:%v,\tnode=`%v`,\t%+v\n", i, pth, pth[i], node, recursive.Pair)
		if node == nil || i == len(pth)-1 {
			inlog.Debugf("- GetByPth > recursive(final),\tpth=%v:%v,\tnode=`%v`,\t%+v\n", pth, pth[i], node, recursive.Pair)
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

// IsSelfFiled returns whether the Col is in the range of the ColHeader
// and does not overlap with the children
// if considerChildrenMarkAsParentField is true, the mark will be counted
// as a column of the ColHeader, otherwise it will be ignored as a child
// and the column will be counted
//
// e.g.:
//
//	ColHeader: [0, 4], Children: [1, 4], IsSelfFiled(1, true): 1 => true
//	ColHeader: [0, 4], Children: [1, 4], IsSelfFiled(1, false): 1 => false
func (header *ColHeader) IsSelfFiled(col Col, considerChildrenMarkAsParentField bool) bool {
	// in range of header.Pair
	if header.Pair != nil && !header.Pair.Inside(col) {
		return false
	}

	contains := func(col Col, pair *pmark.Pair[Col]) bool {
		// overlap with children
		if considerChildrenMarkAsParentField {
			// the left mark is not counted as a part of the child
			return pair.Between(col, false, true)
		} else {
			// the left mark is counted as a part of the child
			return pair.Inside(col)
		}
	}

	// not overlap with children
	for _, child := range header.Children {
		if contains(col, child.Pair) {
			return false
		}
	}
	return true
}

func (header *ColHeader) ForeachCol(action func(colMeta *ColMeta), includeChildren bool) {
	for _, def := range header.Def {
		if !header.IsSelfFiled(def.Col, false) {
			continue
		}
		action(def)
	}
	if !includeChildren {
		return
	}
	// every child can be scanned
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
