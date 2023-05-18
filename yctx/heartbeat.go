package yctx

import (
	"context"
	"mkuznets.com/go/ytils/ytime"
	"time"
)

const (
	DefaultCheckInterval = 10 * time.Second
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

	go func(last *time.Time) {
		for {
			if h.ctx.Err() != nil {
				return
			}

			idle := time.Since(*last)
			if idle >= h.timeout {
				h.cancel()
				return
			}

			left := h.timeout - idle
			sleep := h.checkInterval
			if left < h.checkInterval {
				sleep = left
			}
			ytime.Sleep(h.ctx, sleep)
		}
	}(&lastBeat)

	go func() {
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
