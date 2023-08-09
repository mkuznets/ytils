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
	outputFile := os.Stderr

	hopts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	var shandler slog.Handler
	if fd := outputFile.Fd(); isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
		shandler = slog.NewTextHandler(outputFile, hopts)
	} else {
		shandler = slog.NewJSONHandler(outputFile, hopts)
	}
	slogger := slog.New(shandler)
	slog.SetDefault(slogger)
}

func Err(err error) slog.Attr {
	return slog.Any("err", err)
}
