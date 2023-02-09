package ylog

import (
	"context"
	"github.com/mattn/go-isatty"
	"golang.org/x/exp/slog"
	"os"
)

type ctxKey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// Ctx returns the Logger associated with the ctx
func Ctx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func Setup() {
	var shandler slog.Handler
	outputFile := os.Stderr

	hopts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	if isatty.IsTerminal(outputFile.Fd()) || isatty.IsCygwinTerminal(outputFile.Fd()) {
		shandler = hopts.NewTextHandler(outputFile)
	} else {
		shandler = hopts.NewJSONHandler(outputFile)
	}
	slogger := slog.New(shandler)
	slog.SetDefault(slogger)
}
