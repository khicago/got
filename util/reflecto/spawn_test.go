package reflecto_test

import (
	"testing"

	"github.com/khicago/got/util/reflecto"

	"github.com/stretchr/testify/assert"
)

func TestNewSpawner(t *testing.T) {
	type ty struct {
		N int
	}
	s := ty{}
	sp := reflecto.NewAnySpawner(s)
	v := sp.Spawn()
	x, ok := v.(*ty)
	assert.True(t, ok, "type error")
	assert.Equal(t, s, *x, "spawn error")
}

func TestNewSpawnerPlain(t *testing.T) {
	s := 1
	sp := reflecto.NewAnySpawner(s)
	v := sp.Spawn()
	x, ok := v.(*int)
	assert.True(t, ok, "type error")
	assert.Equal(t, 0, *x, "spawn error")
}
