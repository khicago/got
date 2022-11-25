package pmark_test

// test methods of mark with testify

import (
	"testing"

	"github.com/khicago/got/table2d/preset/pmark"
	"github.com/stretchr/testify/assert"
)

func TestMark(t *testing.T) {
	// ""
	// "<" and ">"
	// "(" and ")"
	// "[" and "]"
	// "{" and "}"
	testDataRegistered := []struct {
		mark pmark.Mark
		pair pmark.Mark
	}{
		{mark: pmark.Mark("<"), pair: pmark.Mark(">")},
		{mark: pmark.Mark("("), pair: pmark.Mark(")")},
		{mark: pmark.Mark("["), pair: pmark.Mark("]")},
		{mark: pmark.Mark("{"), pair: pmark.Mark("}")},
	}

	testDataNotRegistered := []struct {
		mark pmark.Mark
	}{
		{mark: pmark.Mark("")},
		{mark: pmark.Mark("sd")},
	}

	// test testDataRegistered
	for _, td := range testDataRegistered {
		assert.True(t, td.mark.Registered(), td)
		assert.True(t, td.pair.Registered(), td)
		assert.True(t, td.mark.PairedWith(td.pair), td)
		assert.True(t, td.pair.PairedWith(td.mark), td)
		assert.True(t, td.mark.IsLeft(), td)
		assert.False(t, td.pair.IsLeft(), td)
		// IsPair
		assert.True(t, pmark.IsPair(td.mark, td.pair), td)
	}

	// test testDataNotRegistered
	for _, td := range testDataNotRegistered {
		assert.False(t, td.mark.Registered())
		assert.False(t, td.mark.PairedWith(""))
		assert.False(t, td.mark.IsLeft())
	}
}
