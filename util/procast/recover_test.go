package procast_test

import (
	"testing"

	"github.com/khicago/got/util/procast"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	err := func() (err error) {
		defer procast.GetRewriteErrHandler(&err).Recover()
		func() {
			panic("recover, ok, x")
		}()
		return
	}()
	assert.NotNil(t, err, "recover failed")
	assert.Equal(t, "recover, ok, x", err.Error(), "recover content failed")
}
