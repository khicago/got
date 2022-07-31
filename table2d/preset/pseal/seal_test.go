package pseal_test

import (
	"testing"

	"github.com/khicago/got/table2d/preset/pseal"
	"github.com/stretchr/testify/assert"
)

func TestSeal(t *testing.T) {
	v, err := pseal.PID(100010).PID()
	assert.Nil(t, err, "convert pid failed")
	assert.Equal(t, int64(100010), v, "convert pid error")

	v, err = pseal.ID(200300).ID()
	assert.Nil(t, err, "convert id failed")
	assert.Equal(t, int64(200300), v, "convert id error")

	vi, err := pseal.Int(123).Int()
	assert.Nil(t, err, "convert int failed")
	assert.Equal(t, 123, vi, "convert int error")

	vf, err := pseal.Float(1.23).Float()
	assert.Nil(t, err, "convert float failed")
	assert.Equal(t, 1.23, vf, "convert float error")
}

func TestSealByString(t *testing.T) {
	seal, err := pseal.TyPID.SealStr("100010")
	assert.Nil(t, err, "seal pid by string failed")
	v, err2 := seal.PID()

	assert.Nil(t, err2, "seal pid by string, fetching failed")
	assert.Equal(t, int64(100010), v, "convert pid error")
}
