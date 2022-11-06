package basealphabet

import (
	"bytes"

	"github.com/khicago/got/util/typer"
)

func EncodeInt64(alphabet []byte, num int64) string {
	result := make([]byte, 0)
	length := int64(len(alphabet))
	for mode := 0; num != 0; num = num / length {
		mode = int(num % length)
		result = append(result, alphabet[mode])
	}
	typer.SliceReverse(result)
	ret := string(result)
	// inlog.Debugf("ret= %s\n", ret)
	if ret == "" {
		return string(alphabet[:1])
	}
	return ret
}

func DecodeInt64(alphabet []byte, data string) int64 {
	ret := int64(0)
	for _, b := range []byte(data) {
		ret = ret*int64(len(alphabet)) + int64(bytes.IndexByte(alphabet, b))
	}
	return ret
}
