package pcol

import "errors"

var (
	ErrColHeaderUnMarshalFmtFailed = errors.New("col header marshal format failed")
	ErrColMetaUnMarshalFmtFailed   = errors.New("col meta unmarshal format failed")
)
