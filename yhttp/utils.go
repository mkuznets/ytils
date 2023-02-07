package yhttp

import (
	"golang.org/x/exp/slog"
	"net/http"
)

func CloseResponseBody(resp *http.Response) {
	if resp == nil {
		return
	}
	if err := resp.Body.Close(); err != nil {
		slog.Warn("close response body", "err", err)
	}
}
