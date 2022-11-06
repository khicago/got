package transkey

import (
	"context"
	"fmt"

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

func (kl KeyList) ExtractCtxToMap(ctx context.Context, u Usage, marshal delegate.Convert[any, string]) map[string]string {
	ret := make(map[string]string)
	kl.ExtractFromCtx(ctx, u, func(keyForUsage string, value any) {
		if str := marshal(value); str != "" {
			ret[keyForUsage] = str
		}
	})
	return ret
}

func (kl KeyList) ExtractCtxMapStrings(ctx context.Context, u Usage) map[string]string {
	return kl.ExtractCtxToMap(ctx, u, func(val any) string {
		if val == nil {
			return ""
		} else if str, ok := val.(string); ok {
			return str
		} else if stringer, ok := val.(fmt.Stringer); ok {
			return stringer.String()
		}
		return fmt.Sprintf("%v", val)
	})
}

// InjectMapToCtx
// Keys that do not exist in the incoming map are not written
func (kl KeyList) InjectMapToCtx(ctx context.Context, u Usage, m map[string]string, unmarshal func(keyInMap string, v string) any) context.Context {
	ctx = kl.InjectToCtx(ctx, u, func(keyForUsage string) any {
		v, ok := m[keyForUsage]
		if ok {
			return unmarshal(keyForUsage, v)
		}
		return nil
	})
	return ctx
}
