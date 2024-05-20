package idgen

import (
	"context"
	"time"

	"github.com/khicago/got/util/syncounter"
)

type (
	LocalMUGen struct {
		syncounter.Counter
		scopeID        int64
		lastTimePrefix int64
		globalMode     bool // 标记 control 位为 1, 支持全球 1 << 13 个实例
	}
)

// NewLocalMUGen
// fallback: 32b timestamp(second) + 10b timestamp(millisecond) + 8b counter + 00 + 12b scope_id
// 分为运行在 global server 和运行在业务 server 两种模式, 后者为 fallback
// todo: 允许 counter 向 timestamp 借位
func NewLocalMUGen(scopeID int64, fallback bool) *LocalMUGen {
	return &LocalMUGen{
		scopeID:    scopeID,
		globalMode: !fallback,
		Counter:    syncounter.MakeCounter(1<<CounterDigits - 1),
	}
}

func (gen *LocalMUGen) Get(ctx context.Context) (int64, error) {
	timeShifted, idCombine := gen.updatePrefix()
	count, err := gen.CountOne(func() bool { return timeShifted })
	if err != nil {
		return 0, err
	}

	return idCombine | count<<ControlIDDigits, nil
}

func (gen *LocalMUGen) MGet(ctx context.Context, count int64) ([]int64, error) {
	timeShifted, idCombine := gen.updatePrefix()
	form, to, err := gen.Count(count, func() bool { return timeShifted })
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, count)
	for i := form; i <= to; i++ {
		ret = append(ret, idCombine|i<<(ControlIDDigits))
	}
	return ret, nil
}

func (gen *LocalMUGen) updatePrefix() (timeShifted bool, idCombine int64) {
	idCombine = getIDPrefix(time.Now(), CounterDigits)
	lastTimePrefix := gen.lastTimePrefix
	gen.lastTimePrefix = idCombine

	timeShifted = idCombine > lastTimePrefix
	idCombine |= gen.scopeID
	if gen.globalMode {
		idCombine |= ControlDigitsPos
	}
	return timeShifted, idCombine
}
