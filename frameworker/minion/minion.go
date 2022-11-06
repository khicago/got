package minion

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/procast"
)

type (
	InitStrategy byte

	Minion interface {
		Name() string
		ID() string
		Close()
		IsFinished() bool
	}

	minion struct {
		id       string
		name     string
		finished bool

		conf

		chClose   chan struct{}
		watchOnce *sync.Once
	}
)

var _ Minion = &minion{}

const (
	InitAtLeastOnce InitStrategy = 0
	InitAsync       InitStrategy = 1
)

var DefaultConf = conf{
	TickDuration: time.Second,
	InitStrategy: InitAtLeastOnce,
	ErrorHandler: func(err error) {},
	ProcPrinter:  func(w Minion, str string, round int64) {},
}

func New(name string, fnWorkload delegate.Handler, opts ...Option) Minion {
	c := DefaultConf
	for _, opt := range opts {
		c = opt(c)
	}
	minion := &minion{
		conf:      c,
		chClose:   make(chan struct{}),
		watchOnce: &sync.Once{},
	}
	minion.name = name
	rand.Seed(time.Now().UnixNano())

	minion.id = getIDStr(context.Background())
	minion.watchOnce.Do(delegate.WrapAction(minion.start).Partial(fnWorkload))
	return minion
}

func (w *minion) Name() string {
	return w.name
}

func (w *minion) ID() string {
	return w.id
}

func (w *minion) IsFinished() bool {
	return w.finished
}

func (w *minion) Close() {
	close(w.chClose)
}

func (w *minion) start(fnWorkload delegate.Handler) {
	w.ProcPrinter(w, "minion start", 0)
	// once exit called, the holding will release
	// but the goroutine will go on until close
	_ = procast.HoldGo(func(exitHold procast.ErrorHandler) {
		count := int64(0)

		switch w.InitStrategy {
		case InitAsync:
			exitHold(nil)
		case InitAtLeastOnce:
			defer exitHold(nil)
		}

		defer procast.Recover(w.ErrorHandler)
		defer w.ProcPrinter(w, "minion exit", count)

		holdAndTickUntilClose(w.TickDuration, func() {
			count++
			w.ProcPrinter(w, "minion run", count)
			defer procast.Recover(w.ErrorHandler)
			if err := fnWorkload(); err != nil && w.ErrorHandler != nil { // err will not lead the closing of ticks
				w.ErrorHandler(err)
			}
			exitHold(nil) // for the first return (a.k.a. atLeastOnce)
		}, w.chClose)

		w.finished = true
	})
}
