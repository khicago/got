package graceful

import (
	"context"
	"os"
	"os/signal"

	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/procast"
)

func GracefulServe(ctx context.Context,
	blockAndServe delegate.Handler,
	stop delegate.Action,
	errorHandler procast.ErrorHandler,
	signals ...os.Signal,
) error {
	ctx, cancel := WithWaitSignalAndCancel(ctx, signals...)
	defer cancel()

	procast.SafeGo(func() {
		<-ctx.Done()
		stop.TryCall()
	}, errorHandler)
	return blockAndServe()
}

func WithWaitSignal(ctx context.Context, signals ...os.Signal) context.Context {
	ctx, _ = WithWaitSignalAndCancel(ctx, signals...)
	return ctx
}

func WithWaitSignalAndCancel(ctx context.Context, signals ...os.Signal) (context.Context, delegate.Action) {
	ctx, cancel := context.WithCancel(ctx)
	procast.SafeGo(func() {
		WaitSignal(signals...)
		cancel()
	}, nil)
	return ctx, delegate.Action(cancel)
}

func WaitSignal(signals ...os.Signal) os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	return <-c
}
