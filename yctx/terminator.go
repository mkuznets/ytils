package yctx

import (
	"context"
	"github.com/dlsniper/debugger"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

func WithTerminator(ctx context.Context) (context.Context, context.Context) {
	ctx, normalCancel := context.WithCancel(ctx)
	critCtx, criticalCancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 4)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		debugger.SetLabels(func() []string {
			return []string{
				"pkg", "ytils/ycli",
				"func", "context terminator",
			}
		})

		cnt := 0
		for s := range signalChan {
			switch cnt {
			case 0:
				slog.Debug("graceful exit", "signal", s.String())
				normalCancel()
			case 1:
				slog.Debug("send one more for hard exit", "signal", s.String())
				criticalCancel()
			default:
				slog.Debug("hard exit requested, exiting", "signal", s.String())
				os.Exit(1)
			}
			cnt++
		}
	}()

	return ctx, critCtx
}
