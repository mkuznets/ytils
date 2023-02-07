package ycli

import (
	"github.com/jessevdk/go-flags"
	"golang.org/x/exp/slog"
	"os"
)

type Command[T any] interface {
	Init(app *T) error
	Validate() error
	Execute(args []string) error
}

func Main[T any]() {
	var app T
	var parser = flags.NewParser(&app, flags.Default)
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		c := command.(Command[T])
		if err := c.Validate(); err != nil {
			slog.Error("invalid arguments", err)
			os.Exit(1)
		}

		if err := c.Init(&app); err != nil {
			slog.Error("init", err)
			os.Exit(1)
		}

		if err := c.Execute(args); err != nil {
			slog.Error("command error", err)
			os.Exit(1)
		}

		return nil
	}

	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
		slog.Error("invalid arguments", err)
		os.Exit(1)
	}
}
