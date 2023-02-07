package yexec

import (
	"bufio"
	"context"
	"github.com/dlsniper/debugger"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
	"io"
	"os/exec"
	"syscall"
	"time"
)

type signalDelay struct {
	signal syscall.Signal
	delay  time.Duration
}

var terminationSequence = []signalDelay{
	{syscall.SIGINT, 10 * time.Second},
	{syscall.SIGTERM, 30 * time.Second},
	{syscall.SIGKILL, 0},
}

type Execer struct {
	hardTimeout   *time.Duration
	stdoutFunc    func(string)
	stderrFunc    func(string)
	sigintContext context.Context
}

func NewExecer() *Execer {
	ex := &Execer{}
	return ex
}

func (e *Execer) WithStdoutFunc(f func(string)) *Execer {
	e.stdoutFunc = f
	return e
}

func (e *Execer) WithStderrFunc(f func(string)) *Execer {
	e.stderrFunc = f
	return e
}

func (e *Execer) WithGracefulExit(ctx context.Context) *Execer {
	e.sigintContext = ctx
	return e
}

func (e *Execer) Exec(cmd *exec.Cmd) error {
	var (
		stdout, stderr io.ReadCloser
		err            error
	)

	if e.stdoutFunc != nil {
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}
	}
	if e.stderrFunc != nil {
		stderr, err = cmd.StderrPipe()
		if err != nil {
			return err
		}
	}

	g := new(errgroup.Group)

	if err := cmd.Start(); err != nil {
		return err
	}

	pLogger := slog.With("pid", cmd.Process.Pid)

	if stdout != nil {
		g.Go(func() error {
			debugger.SetLabels(func() []string {
				return []string{
					"pkg", "ytils/yexec",
					"cmd", cmd.String(),
					"func", "stdout reader",
				}
			})
			readLines(stdout, e.stdoutFunc)
			return nil
		})
	}

	if stderr != nil {
		g.Go(func() error {
			debugger.SetLabels(func() []string {
				return []string{
					"pkg", "ytils/yexec",
					"cmd", cmd.String(),
					"func", "stderr reader",
				}
			})
			readLines(stderr, e.stderrFunc)
			return nil
		})
	}

	sic, cancel := context.WithCancel(e.sigintContext)
	defer cancel()
	go func() {
		debugger.SetLabels(func() []string {
			return []string{
				"pkg", "ytils/yexec",
				"cmd", cmd.String(),
				"func", "sigintContext monitor",
			}
		})

		select {
		case <-e.sigintContext.Done():
			for _, s := range terminationSequence {
				if cmd.ProcessState != nil {
					// Process already exited
					return
				}
				if err := cmd.Process.Signal(s.signal); err != nil {
					pLogger.Error(
						"failed to send signal",
						err,
						"signal", s.signal.String(),
						"error", err,
					)
				} else {
					pLogger.Debug("signal sent", "signal", s.signal.String())
				}

				pLogger.Debug("waiting before next signal", "delay", s.delay)
				time.Sleep(s.delay)
			}

			return
		case <-sic.Done():
			return
		}
	}()

	// Wait for the output readers to finish before exiting since cmd.Wait() closes the readers.
	if err := g.Wait(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func readLines(r io.Reader, callback func(string)) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		callback(sc.Text())
	}
}
