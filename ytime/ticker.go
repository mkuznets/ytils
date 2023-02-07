package ytime

import (
	"context"
	"time"
)

type Ticker struct {
	skipFirst bool
	interval  time.Duration
}

func NewTicker(interval time.Duration) *Ticker {
	ticker := &Ticker{
		skipFirst: false,
		interval:  interval,
	}
	return ticker
}

func (t *Ticker) SkipFirst() *Ticker {
	t.skipFirst = true
	return t
}

func (t *Ticker) Start(ctx context.Context, fun func() error) error {
	if !t.skipFirst {
		if err := fun(); err != nil {
			return err
		}
	}

	timer := time.NewTimer(t.interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-timer.C:
			if ctx.Err() != nil {
				timer.Reset(t.interval)
				return nil
			}

			if err := fun(); err != nil {
				return err
			}

			timer.Reset(t.interval)
		}
	}
}
