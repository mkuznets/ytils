package ytime

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return
	case <-timer.C:
		timer.Stop()
		return
	}
}
