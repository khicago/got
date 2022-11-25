package pcol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test methods of pcol with testify
func TestMetaMarshalAndUnMarshal(t *testing.T) {
	meta := NewColMeta(3, "ID", "test", ">0")

	// MarshalJSON
	b, err := meta.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "3:id:test:>0", string(b))

	// UnmarshalJSON
	var meta2 ColMeta
	err = meta2.UnmarshalJSON(b)
	assert.NoError(t, err)
	assert.Equal(t, *meta, meta2)
}
