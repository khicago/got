package preset

import (
	"testing"

	"github.com/khicago/got/table2d/preset/pcol"
	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/khicago/got/util/inlog"

	"github.com/khicago/got/table2d/preset/pseal"

	"github.com/stretchr/testify/assert"
)

func TestPreset(t *testing.T) {
	example := &Prop{
		p: PropData{
			0:  pseal.PID(100010), // pid
			1:  pseal.ID(200300),  // id
			2:  pseal.Int(123),    // int
			3:  pseal.Float(1.23),
			4:  pseal.Mark("["),
			5:  pseal.Int(1),
			6:  pseal.Int(2),
			7:  pseal.Int(3),
			8:  pseal.Int(4),
			9:  pseal.Mark("]"),
			10: pseal.ID(200303), // id
		},
		childrenCols: PropChildIndex{
			4: pmark.Pair[pcol.Col]{
				L:    "{",
				R:    "}",
				LVal: 4,
				RVal: 9,
			},
		},
	}

	v, err := example.Get(0).PID()
	assert.Nil(t, err, "convert pid failed")
	assert.Equal(t, int64(100010), v, "convert pid error")

	v, err = example.Get(1).ID()
	assert.Nil(t, err, "convert id failed")
	assert.Equal(t, int64(200300), v, "convert id error")

	vi, err := example.Get(2).Int()
	assert.Nil(t, err, "convert int failed")
	assert.Equal(t, 123, vi, "convert int error")

	vf, err := example.Get(3).Float()
	assert.Nil(t, err, "convert float failed")
	assert.Equal(t, 1.23, vf, "convert float error")

	vList, err := example.Child(4)
	assert.Nil(t, err, "convert list failed")

	testValInd := 0
	vList.ForEach(func(col pcol.Col, seal pseal.Seal) {
		testValInd++
		vListVal, err := seal.Int()
		assert.Nil(t, err, "convert list failed")
		assert.Equal(t, testValInd, vListVal, "convert list child val error")

		inlog.Infof("list val %v: %v %v", testValInd, vListVal, err)
	})
	//
	//voi, err := vl.Get(5).Int()
	//assert.Equal(t, 1, voi, "convert child param error")
	//
	//vid, err := example.Get(5).PID()
	//assert.Nil(t, err, "convert id fallback failed")
	//assert.Equal(t, int64(200303), vid, "convert id fallback error")
}
