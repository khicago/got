package presetor

import (
	"context"
	"github.com/khicago/got/table2d/preset"
	"github.com/khicago/got/table2d/tablety"
)

func Strings2d(ctx context.Context, table [][]string) (*preset.Preset, error) {
	return preset.ReadLines(ctx, tablety.WarpLineReader(table))
}
