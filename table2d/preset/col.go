package preset

import "github.com/khicago/got/table2d/preset/pseal"

type (
	ColMeta struct {
		Col
		Name       string
		Type       pseal.Type
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
