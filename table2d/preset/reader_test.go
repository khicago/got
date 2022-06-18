package preset

import (
	"context"
	"fmt"
	"github.com/khicago/got/internal/utils"
	"github.com/khicago/got/table2d/tablety"
	"github.com/khicago/got/util/typer"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	MockTableReader struct {
		Data [][]string
	}
)

func (m MockTableReader) MaxRow() int {
	return len(m.Data) - 1
}

func (m MockTableReader) MaxCol() int {
	if len(m.Data) == 0 {
		return 0
	}
	return len(m.Data[0]) - 1
}

func (m MockTableReader) Get(row tablety.Row, col tablety.Col) string {
	return m.Data[row][col]
}

func (m MockTableReader) First(pred typer.Predicate[string]) (tablety.Row, tablety.Col) {
	for r := 0; r <= m.MaxRow(); r++ {
		for c := 0; c <= m.MaxRow(); c++ {
			if pred(m.Get(r, c)) {
				return r, c
			}
		}
	}
	return -1, -1
}

var _ tablety.Table2DReader[string] = &MockTableReader{}

var data = MockTableReader{
	Data: [][]string{
		{"@", "ID", "INT", "Float", "[", "", "]"},
		{" ", "link(@)", "test($>1,$<50)", "test($%2)", "link(item)", "", ""},
		{"", "LvUp", "Power", "Magic", "InitItems", "", ""},
		{"10001", "10002", "12", "1.2", "", "1010001", "1010002"},
	},
}

func TestPresetReader(t *testing.T) {
	p := Read(context.TODO(), data)
	fmt.Printf("header: %s\n", utils.MarshalPrintAll(p.Header))
	fmt.Printf("table: %s\n", utils.MarshalPrintAll(p.PropTable))

	lvUpMeta := p.Header.GetByName("lv_up")
	if !assert.NotNil(t, lvUpMeta, "lv_up col cannot found") {
		return
	}

	lvUpCol := lvUpMeta.Col
	if !assert.Equal(t, 1, lvUpCol, "lv_up col error") {
		return
	}
	v, err := p.Query(10001, lvUpCol).ID()
	if !assert.Nil(t, err, "get lvUpCol of 10001 failed") {
		return
	}
	assert.Equal(t, int64(10002), v, "get lvUpCol of 10001 val error")
}
