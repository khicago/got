package idgen

import (
	"context"
	"github.com/khicago/got/util/syncounter"
	"time"
)

type (
	FakeGen struct {
		syncounter.Counter
		lastTimePrefix int64
	}
)

var _ IGenerator = NewFakeGen()

func NewFakeGen() *FakeGen {
	return &FakeGen{
		Counter: syncounter.MakeCounter(1<<CounterDigits - 1),
	}
}

func (gen *FakeGen) Get(ctx context.Context) (int64, error) {
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

func (gen *FakeGen) MGet(ctx context.Context, count int64) ([]int64, error) {
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
