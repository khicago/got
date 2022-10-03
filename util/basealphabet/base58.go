package basealphabet

import "github.com/khicago/got/util/delegate"

var (
	Base58BitCoin        = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	Base58BitCoinEncoder = delegate.Func2[[]byte, int64, string](EncodeInt64).Partial(Base58BitCoin)
	Base58BitCoinDecoder = delegate.Func2[[]byte, string, int64](DecodeInt64).Partial(Base58BitCoin)
)
