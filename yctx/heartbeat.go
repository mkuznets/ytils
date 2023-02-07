package yctx

import (
	"context"
	"github.com/dlsniper/debugger"
	"golang.org/x/exp/slog"
	"mkuznets.com/go/ytils/ylog"
	"mkuznets.com/go/ytils/ytime"
	"time"
)

const (
	DefaultCheckInterval = 10 * time.Second
	DefaultLeftWarning   = time.Minute
)

type Heartbeat struct {
	beatC         chan bool
	ctx           context.Context
	cancel        context.CancelFunc
	timeout       time.Duration
	checkInterval time.Duration
	leftWarning   time.Duration
}

func NewHeartbeat(ctx context.Context, timeout time.Duration) *Heartbeat {
	ctx, cancel := context.WithCancel(ctx)
	return &Heartbeat{
		ctx:           ctx,
		cancel:        cancel,
		timeout:       timeout,
		checkInterval: DefaultCheckInterval,
		leftWarning:   DefaultLeftWarning,
		beatC:         make(chan bool),
	}
}

func (h *Heartbeat) Context() context.Context {
	return h.ctx
}

func (h *Heartbeat) Beat() {
	select {
	case h.beatC <- true:
	default:
	}
}

func (h *Heartbeat) Close() {
	h.cancel()
	close(h.beatC)
}

func (h *Heartbeat) Start() *Heartbeat {
	lastBeat := time.Now()
	logger := ylog.Ctx(h.ctx)

	go func(last *time.Time) {
		debugger.SetLabels(func() []string {
			return []string{
				"pkg", "ytils/yctx",
				"func", "beats monitor",
			}
		})

		for {
			if h.ctx.Err() != nil {
				return
			}

			idle := time.Since(*last)
			slog.Debug("heartbeat", "elaplsed", idle)

			if idle >= h.timeout {
				h.cancel()
				logger.Warn("idle context cancelled")
				return
			}

			left := h.timeout - idle
			if left <= DefaultLeftWarning {
				logger.Warn("idle context", "left", left)
			}

			sleep := h.checkInterval
			if left < h.checkInterval {
				sleep = left
			}
			ytime.Sleep(h.ctx, sleep)
		}
	}(&lastBeat)

	go func() {
		debugger.SetLabels(func() []string {
			return []string{
				"pkg", "ytils/yctx",
				"func", "beats consumer",
			}
		})

		for {
			select {
			case <-h.beatC:
				lastBeat = time.Now()
			case <-h.ctx.Done():
				return
			}
		}
	}()

	return h
}
