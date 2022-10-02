package idgen

import (
	"context"
	"errors"
)

type (
	IGenerator interface {
		Get(ctx context.Context) (int64, error)
		MGet(ctx context.Context, count int64) ([]int64, error)
	}
)

// (集中 server 部署情况下, 不需要区分 scope)
// normal: 32b timestamp(second) + 10b timestamp(millisecond) + 8b counter + 14b server_id (1 + 13b id_server_id)
// custom: 32b timestamp(second) + 18b custom_counter + 01 + 12b scope_id
// fallback: 32b timestamp(second) + 10b timestamp(millisecond) + 8b counter + 00 + 12b scope_id
// fake: 32b timestamp(second) + 10b timestamp(millisecond) + 8b counter + 00 + 000000000000

const (
	ControlIDDigits = 14
	CounterDigits   = 8

	ControlDigitsPos = 1 << 13
)

var (
	ErrGenerateIDFailed = errors.New("generate id failed")
)
