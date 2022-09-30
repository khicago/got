package preset

import (
	"context"
	"github.com/khicago/got/internal/utils"
	"github.com/khicago/got/table2d/tablety"
	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/inlog"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	MockTableReader struct {
		Data [][]string
	}

	MockLineReader struct {
		*MockTableReader
		row int
	}
)

func (m MockTableReader) LineReader() tablety.LineReader[string] {
	return &MockLineReader{MockTableReader: &m}
}

func (m *MockLineReader) Read() (ret []string, err error) {
	if m.row < len(m.MockTableReader.Data) {
		ret = m.MockTableReader.Data[m.row]
		m.row++
	}
	return ret, nil
}

func (m MockTableReader) Reader() *MockLineReader {
	return &MockLineReader{MockTableReader: &m}
}

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

func (m MockTableReader) First(pred delegate.Predicate[string]) (tablety.Row, tablety.Col) {
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
		{"@", "ID", "INT", "Float", "[", "ID", "]", "{", "ID", "}"},
		{" ", "link(@)", "test($>1,$<50)", "test($%2)", "link(item)", "", "", "select", "", ""},
		{"", "LvUp", "Power", "Magic", "InitItems", "", "", "InnerLvUpItem", "LvUp", ""},
		{"10001", "10002", "12", "1.2", "", "1010001", "", "", "1010003", ""},
	},
}

func TestPresetReader(t *testing.T) {
	p, err := Read(context.TODO(), data)
	if !assert.Nil(t, err, "read by table failed") {
		return
	}

	inlog.Infof("header: %s\n", utils.MarshalIndentPrintAll(p.Headline))
	inlog.Infof("table: %s\n", utils.MarshalPrintAll(p.PropTable))

	lvUpMeta := p.Headline.GetByPth("lv_up")
	if !assert.NotNil(t, lvUpMeta, "lv_up col cannot found") {
		return
	}

	lvUpCol := lvUpMeta.Col
	if !assert.Equal(t, 1, lvUpCol, "lv_up col error") {
		return
	}
	v, err := p.QueryByCol(10001, lvUpCol).ID()
	if !assert.Nil(t, err, "get lvUpCol of 10001 failed") {
		return
	}
	assert.Equal(t, int64(10002), v, "get lvUpCol of 10001 val error")
}

func TestPresetReaderLines(t *testing.T) {
	p, err := ReadLines(context.TODO(), data.LineReader())
	if !assert.Nil(t, err, "read by lineReader failed") {
		return
	}
	inlog.Infof("header: %s\n", utils.MarshalIndentPrintAll(p.Headline))
	inlog.Infof("table: %s\n", utils.MarshalPrintAll(p.PropTable))

	lvUpMeta := p.Headline.GetByPth("lv_up")
	if !assert.NotNil(t, lvUpMeta, "lv_up col cannot found") {
		return
	}

	lvUpCol := lvUpMeta.Col
	if !assert.Equal(t, 1, lvUpCol, "lv_up col error") {
		return
	}
	v, err := p.QueryByCol(10001, lvUpCol).ID()
	if !assert.Nil(t, err, "get lvUpCol of 10001 failed") {
		return
	}
	assert.Equal(t, int64(10002), v, "get lvUpCol of 10001 val error")

	InitItems0 := p.Query(10001, "init_items", "0")
	if !assert.NotNil(t, InitItems0, "init_items col cannot found") {
		return
	}
	v, err = InitItems0.ID()
	if !assert.Nil(t, err, "get init_items[0] of 10001 failed") {
		return
	}
	assert.Equal(t, int64(1010001), v, "get init_items[0] of 10001 val error")

	InnerLvUpItemID := p.QueryS(10001, "inner_lv_up_item/lv_up")
	if !assert.NotNil(t, InnerLvUpItemID, "inner_lv_up_item/lv_up col cannot found") {
		return
	}
	v, err = InnerLvUpItemID.ID()
	if !assert.Nil(t, err, "get inner_lv_up_item/lv_up of 10001 failed") {
		return
	}
	assert.Equal(t, int64(1010003), v, "get inner_lv_up_item/lv_up of 10001 val error")
}
