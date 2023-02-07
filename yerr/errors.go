package yerr

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"runtime"
	"strings"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type Error interface {
	error
	StackTracer
	Status() int
	Err(error) Error
	Unwrap() error
	Message() string
}

type errorImpl struct {
	err     error
	status  int
	message string
	stack   []uintptr
}

func (e *errorImpl) Unwrap() error {
	return e.err
}

func (e *errorImpl) Err(err error) Error {
	e.err = err
	return e
}

func (e *errorImpl) Status() int {
	return e.status
}

func (e *errorImpl) Message() string {
	return e.message
}

func (e *errorImpl) Error() string {
	var (
		msgs []string
		err  error = e
	)
	for err != nil {
		switch v := err.(type) {
		case Error:
			msgs = append(msgs, v.Message())
		default:
			msgs = append(msgs, err.Error())
		}
		err = errors.Unwrap(err)
	}

	return strings.Join(msgs, ": ")
}

func (e *errorImpl) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(e.stack))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame(e.stack[i])
	}
	return f
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(4, pcs[:])
	var st = pcs[0:n]
	return st
}

func newErrorf(err error, status int, message string, a ...interface{}) Error {
	return &errorImpl{err, status, fmt.Sprintf(message, a...), callers()}
}

func NotFound(format string, a ...interface{}) Error {
	return newErrorf(nil, http.StatusNotFound, format, a...)
}

func Invalid(format string, a ...interface{}) Error {
	return newErrorf(nil, http.StatusBadRequest, format, a...)
}

func Unauthorised(format string, a ...interface{}) Error {
	return newErrorf(nil, http.StatusUnauthorized, format, a...)
}

func New(format string, a ...interface{}) Error {
	return newErrorf(nil, http.StatusInternalServerError, format, a...)
}
