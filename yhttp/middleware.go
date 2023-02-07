package yhttp

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/segmentio/ksuid"
	"golang.org/x/exp/slog"
	"mkuznets.com/go/ytils/ylog"
	"net/http"
	"time"
)

type contextKey int

const (
	ctxRequestIdKey = contextKey(0x5245)
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

func RequestId(r *http.Request) string {
	if reqID, ok := r.Context().Value(ctxRequestIdKey).(string); ok {
		return reqID
	}
	return ""
}

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = "req_" + ksuid.New().String()
		}
		ctx = context.WithValue(ctx, ctxRequestIdKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ContextLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.Default()
		if reqId := RequestId(r); reqId != "" {
			logger = logger.With("req_id", reqId)
		}
		ctx := ylog.WithLogger(r.Context(), logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logAttrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.RequestURI),
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			logAttrs = append(logAttrs,
				slog.Duration("duration", time.Since(t1)),
				slog.Int("status", ww.Status()),
				slog.Int("size", ww.BytesWritten()),
			)

			ylog.Ctx(r.Context()).LogAttrs(slog.LevelInfo, "API", logAttrs...)
		}()

		next.ServeHTTP(ww, r)
	})
}
