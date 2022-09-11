package preset

import (
	"encoding/json"
	"fmt"
	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/khicago/got/util/strs"
	"strconv"
	"strings"
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

func (c *ColMeta) String() string {
	return fmt.Sprintf("%d:%s:%s:%v", c.Col, c.Sym, c.Name, c.Constraint)
}

func (c *ColMeta) MarshalJSON() ([]byte, error) {
	str := c.String()
	return json.Marshal(str)
}

func (c *ColMeta) UnmarshalJSON(bytes []byte) error {
	str := string(bytes)
	values := strings.Split(str, ":")[0:]
	if len(values) < 4 {
		return ErrSealFormatError
	}
	col, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	name, sym, constraint := values[1], values[2], str[len(values[0])+len(values[1])+len(values[2])+3:]
	*c = *NewColMeta(col, sym, name, constraint)
	return nil
}
