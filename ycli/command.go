package ycli

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

// PrintErrors is intentionally excluded because it is not applied to internal errors
// (such as default value on a boolean flag).
// We want to print these internal errors, but don't want to print user errors twice.
var flagsOptions flags.Options = flags.HelpFlag | flags.PassDoubleDash

// Initer is an optional interface for commands that require global flags.
// The command can read the global values and store them to be later used in Execute.
type Initer[T any] interface {
	Init(app *T) error
}

// Validator is an optional interface for commands that require additional validation of their flags.
type Validator interface {
	Validate() error
}

func Main[T any]() *T {
	var app T
	var parser = flags.NewParser(&app, flagsOptions)
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if command == nil {
			return nil
		}

		if validator, ok := command.(Validator); ok {
			if err := validator.Validate(); err != nil {
				errorExit(err)
			}
		}

		if initer, ok := command.(Initer[T]); ok {
			if err := initer.Init(&app); err != nil {
				errorExit(err)
			}
		}

		if err := command.Execute(args); err != nil {
			errorExit(err)
		}

		return nil
	}

	if _, err := parser.Parse(); err != nil {
		errorExit(err)
	}

	return &app
}

func errorExit(err error) {
	if err == nil {
		return
	}

	flagsErr, ok := err.(*flags.Error)
	if ok && flagsErr.Type == flags.ErrHelp {
		_, _ = fmt.Fprintln(os.Stdout, err)
		os.Exit(0)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return
}
