package pcol

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/strs"
)

type (
	Col = int

	// ColMeta defines how this column is resolved, only the first Mark
	// of the substructure is recorded in the upper ColMeta list
	ColMeta struct {
		// Col is the column number of the meta, this information is also
		// available in the upper level table, the redundancy is provided
		// here mainly to pass the full information separately and later.
		Col

		// Type is the type of data for the column
		Type pseal.Type

		// Sym used to store the actual symbol filled in, to do some data
		// comparison, and to derive the type. When serialising and
		// de-serialising, the symbol should be serialised and de-serialised
		// in preference to the type.
		Sym string

		// Name column names are automatically converted to snake lower
		// case as the data is tiled, so even if there are substructures,
		// there can be no renamed columns
		Name string

		// Constraint is responsible for data validation, correlation,
		// processing, etc.
		Constraint string
	}
)

const (
	InvalidCol Col = -1
)

// NewColMeta creates a new ColMeta
// sym - the symbol of the column, marks the type (or pair type) of the
// column, and is automatically converted to lower case.
// name - the name of the column, is automatically converted to snake case
// constraint - the constraint of the column, is used to validate the
// data, and is automatically converted to lower case.
func NewColMeta[TCol ~Col](col TCol, sym, name, constraint string) *ColMeta {
	sym = strs.TrimLower(sym)
	name = strs.Conv2Snake(name)
	return &ColMeta{
		Col:        Col(col),
		Sym:        sym,
		Name:       name,
		Type:       pseal.SymToType(sym),
		Constraint: constraint,
	}
}

var (
	_ json.Marshaler   = &ColMeta{}
	_ json.Unmarshaler = &ColMeta{}
	_ fmt.Stringer     = &ColMeta{}
)

// String returns the string representation of the ColMeta
func (c *ColMeta) String() string {
	return fmt.Sprintf("%d:%s:%s:%v", c.Col, c.Sym, c.Name, c.Constraint)
}

func (c *ColMeta) MarshalJSON() ([]byte, error) {
	str := c.String()
	return []byte(str), nil
}

func (c *ColMeta) UnmarshalJSON(bytes []byte) error {
	str := string(bytes)

	values := strings.Split(str, ":")[0:]
	if len(values) < 4 {
		return ErrColMetaUnMarshalFmtFailed
	}
	col, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	sym, name, constraint := values[1], values[2], str[len(values[0])+len(values[1])+len(values[2])+3:]
	*c = *NewColMeta(col, sym, name, constraint)
	return nil
}
