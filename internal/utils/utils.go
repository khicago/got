package utils

import "encoding/json"

func MarshalPrintAll(v any) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}
