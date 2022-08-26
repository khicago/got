package utils

import "encoding/json"

func MarshalPrintAll(v any) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func MarshalIndentPrintAll(v any) string {
	bs, _ := json.MarshalIndent(v, "", "  ")
	return string(bs)
}
