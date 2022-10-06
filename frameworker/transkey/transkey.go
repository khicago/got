package transkey

import (
	"context"
	"github.com/khicago/got/util/delegate"
)

type (
	Usage string

	H map[Usage]string

	KeyGroup struct {
		// runtime
		CtxKey string
		// protocol
		H
	}

	KeyList []*KeyGroup
)

const (
	HEADER Usage = "header"
	RPC    Usage = "rpc"
	LOG    Usage = "log"
)

func NewTransKeyGroup(ctxKey string, usage H) *KeyGroup {
	if usage == nil {
		usage = make(H)
	}
	return &KeyGroup{
		CtxKey: ctxKey,
		H:      usage,
	}
}

func (kg *KeyGroup) KeyOf(u Usage) string {
	return kg.H[u]
}

func (kg *KeyGroup) InjectToCtx(ctx context.Context, val interface{}) context.Context {
	return context.WithValue(ctx, kg.CtxKey, val)
}

func (kg *KeyGroup) ExtractFromCtx(ctx context.Context) interface{} {
	return ctx.Value(kg.CtxKey)
}

func (kl KeyList) InjectToCtx(ctx context.Context, u Usage,
	getter func(keyForUsage string) any,
	fallbacks ...func(keyForUsage string) any,
) context.Context {
	for _, kg := range kl {
		key, ok := kg.H[u]
		if !ok {
			continue
		}
		v := getter(key)
		for i := 0; v == nil && i < len(fallbacks); i++ {
			v = delegate.Func1[string, any](fallbacks[i]).TryCall(key, nil)
		}
		if v != nil {
			ctx = context.WithValue(ctx, kg.CtxKey, v)
		}
	}
	return ctx
}

// ExtractFromCtx
// 找出 ctx 中所有被 KeyList 命中的值, 并返回这些值与对应 usage 的 key
func (kl KeyList) ExtractFromCtx(ctx context.Context, u Usage, setter func(keyForUsage string, value any)) context.Context {
	for _, kg := range kl {
		key, ok := kg.H[u]
		if !ok {
			continue
		}
		setter(key, ctx.Value(kg.CtxKey))
	}
	return ctx
}

func (kl KeyList) ExtractCtxToMap(ctx context.Context, u Usage, marshal delegate.Map[any, string]) map[string]string {
	ret := make(map[string]string)
	kl.ExtractFromCtx(ctx, u, func(keyForUsage string, value any) {
		ret[keyForUsage] = marshal(value)
	})
	return ret
}

func (kl KeyList) InjectMapToCtx(ctx context.Context, u Usage, m map[string]string, unmarshal func(key string, v string) any) context.Context {
	ctx = kl.InjectToCtx(ctx, u, func(keyForUsage string) any {
		return unmarshal(keyForUsage, m[keyForUsage])
	})
	return ctx
}
