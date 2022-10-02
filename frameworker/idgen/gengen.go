package idgen

import (
	"context"
	"fmt"
)

type (
	IDGen struct {
		generators []IGenerator
	}
)

func NewIDGen(generators ...IGenerator) *IDGen {
	if len(generators) == 0 {
		generators = []IGenerator{NewFakeGen()}
	}
	return &IDGen{
		generators: generators,
	}
}

func (g *IDGen) Get(ctx context.Context) (int64, error) {
	err := ErrGenerateIDFailed
	for _, gen := range g.generators {
		v, e := gen.Get(ctx)
		if e != nil {
			err = fmt.Errorf("%w, %s", err, e)
			continue
		}
		return v, nil
	}
	return 0, err
}

func (g *IDGen) MGet(ctx context.Context, count int64) ([]int64, error) {
	err := ErrGenerateIDFailed
	for _, gen := range g.generators {
		ret, e := gen.MGet(ctx, count)
		if e != nil {
			err = fmt.Errorf("%w, %s", err, e)
			continue
		}
		return ret, nil
	}
	return nil, err
}
