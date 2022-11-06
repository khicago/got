package idgen

import (
	"context"
	"time"

	"github.com/khicago/got/util/syncounter"
)

type (
	fakeGen struct {
		syncounter.Counter
		lastTimePrefix int64
	}
)

// NewFakeGen returns a fake generator
// rule: 32b timestamp(second) + 10b timestamp(millisecond) + 8b counter + 00 + 000000000000
func NewFakeGen() IGenerator {
	return &fakeGen{
		Counter: syncounter.MakeCounter(1<<CounterDigits - 1),
	}
}

func (gen *fakeGen) Get(ctx context.Context) (int64, error) {
	prefixTs := getIDPrefix(time.Now(), 0)
	lastTimePrefix := gen.lastTimePrefix
	gen.lastTimePrefix = prefixTs
	count, err := gen.CountOne(func() bool {
		return prefixTs > lastTimePrefix
	})
	if err != nil {
		return 0, err
	}

	prefixTs |= count << ControlIDDigits
	return prefixTs, nil
}

func (gen *fakeGen) MGet(ctx context.Context, count int64) ([]int64, error) {
	prefixTs := getIDPrefix(time.Now(), 0)
	lastTimePrefix := gen.lastTimePrefix
	gen.lastTimePrefix = prefixTs
	form, to, err := gen.Count(count, func() bool {
		return prefixTs > lastTimePrefix
	})
	if err != nil {
		return nil, err
	}

	ret := make([]int64, count)
	for i := form; i <= to; i++ {
		ret = append(ret, prefixTs|(i<<ControlIDDigits))
	}
	return ret, nil
}
