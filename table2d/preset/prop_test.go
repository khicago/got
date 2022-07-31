package preset

import (
	"testing"

	"github.com/khicago/got/table2d/preset/pseal"

	"github.com/stretchr/testify/assert"
)

func TestPreset(t *testing.T) {
	example := &Prop{
		0: pseal.PID(100010), // pid
		1: pseal.ID(200300),  // id
		2: pseal.Int(123),    // int
		3: pseal.Float(1.23),
		4: pseal.Any(&Prop{ // list
			0: pseal.Any(&Prop{ // obj
				0: pseal.Int(1),
			}),
		}),
		5: pseal.ID(200303), // id
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

	vl, err := example.Child(4)
	assert.Nil(t, err, "convert list failed")

	vo, err := vl.Child(0)
	assert.Nil(t, err, "convert object failed")

	voi, err := vo.Get(0).Int()
	assert.Equal(t, 1, voi, "convert child param error")

	vid, err := example.Get(5).PID()
	assert.Nil(t, err, "convert id fallback failed")
	assert.Equal(t, int64(200303), vid, "convert id fallback error")
}
