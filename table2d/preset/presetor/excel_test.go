package presetor

import (
	"testing"

	"github.com/khicago/got/internal/utils"
	"github.com/khicago/got/util/inlog"
	"github.com/stretchr/testify/assert"
)

func TestExcel(t *testing.T) {
	p, err := ExcelFile("./preset_test.xlsx", "")
	if !assert.Nil(t, err, "read by table failed") {
		return
	}

	inlog.Infof("header: %s\n", utils.MarshalPrintAll(p.Headline))
	inlog.Infof("table: %s\n", utils.MarshalPrintAll(p.PropTable))

	lvUpMeta := p.Headline.GetByPth("lv_up")
	if !assert.NotNil(t, lvUpMeta, "lv_up col cannot found") {
		return
	}

	v, err := p.QueryByCol(10001, lvUpMeta.Col).ID()
	if !assert.Nil(t, err, "get lvUpCol of 10001 failed") {
		return
	}
	assert.Equal(t, int64(10002), v, "get lvUpCol of 10001 val error")

	v, err = p.QueryByCol(10002, lvUpMeta.Col).ID()
	if !assert.Nil(t, err, "get lvUpCol of 10001 failed") {
		return
	}
	assert.Equal(t, int64(10003), v, "get lvUpCol of 10001 val error")
}
