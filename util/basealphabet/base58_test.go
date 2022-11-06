package basealphabet

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase58BitCoinEncode(t *testing.T) {
	str := Base58BitCoinEncoder(0)
	assert.Equal(t, "1", str)
	assert.Equal(t, int64(0), Base58BitCoinDecoder(str))
	str = Base58BitCoinEncoder(1)
	assert.Equal(t, "2", str)
	assert.Equal(t, int64(1), Base58BitCoinDecoder(str))
	str = Base58BitCoinEncoder(99999)
	assert.Equal(t, "Wj8", str)
	assert.Equal(t, int64(99999), Base58BitCoinDecoder(str))
	str = Base58BitCoinEncoder(math.MaxInt64)
	assert.Equal(t, "NQm6nKp8qFC", str)
	assert.Equal(t, int64(math.MaxInt64), Base58BitCoinDecoder(str))
}
