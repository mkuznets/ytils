package ylog

import (
	"fmt"
	"github.com/pkg/errors"
	"mkuznets.com/go/ytils/yerr"
)

var (
	StackSourceLineName     = "src"
	StackSourceFunctionName = "func"
)

// MarshalStack implements pkg/errors stack trace marshaling.
func MarshalStack(err error) interface{} {
	e := err
	var st errors.StackTrace
	for e != nil {
		if errS, ok := e.(yerr.StackTracer); ok {
			st = errS.StackTrace()
		}
		e = errors.Unwrap(e)
	}
	if st == nil {
		return nil
	}

	out := make([]map[string]string, 0, len(st))
	for _, frame := range st {
		out = append(out, map[string]string{
			StackSourceLineName:     fmt.Sprintf("%v", frame),
			StackSourceFunctionName: fmt.Sprintf("%n", frame),
		})
	}
	return out
}
