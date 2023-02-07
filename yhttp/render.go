package yhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"mkuznets.com/go/ytils/yerr"
	"mkuznets.com/go/ytils/ylog"
	"net/http"
	"regexp"
	"strconv"
)

var statusCodeRegex = regexp.MustCompile(`^HTTP (\d{3}):\s*(.+)`)

type Response interface {
	Status(s int) Response
	JSON()
	XML()
}

type response struct {
	w      http.ResponseWriter
	r      *http.Request
	status int
	v      interface{}
}

func Render(w http.ResponseWriter, r *http.Request, v interface{}) Response {
	return &response{
		w:      w,
		r:      r,
		v:      v,
		status: http.StatusOK,
	}
}

func (r *response) Status(status int) Response {
	r.status = status
	return r
}

func (r *response) JSON() {
	switch obj := r.v.(type) {
	case error:
		renderJSONError(r.w, r.r, obj)
	default:
		renderJSON(r.w, r.status, r.v)
	}
}

func (r *response) XML() {
	r.w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	switch obj := r.v.(type) {
	case string:
		r.w.WriteHeader(r.status)
		_, _ = r.w.Write([]byte(obj))
	default:
		http.Error(r.w, "Unsupported type", http.StatusInternalServerError)
	}
}

func renderJSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func reportError(ctx context.Context, err error, stack bool) {
	ylog.Ctx(ctx).Error("internal error", err)
}

func extractStatus(err error) (int, string) {
	msg := err.Error()
	if matches := statusCodeRegex.FindStringSubmatch(msg); matches != nil {
		statusCode := yerr.Must(strconv.Atoi(matches[1]))
		msg = matches[2]
		return statusCode, msg
	}
	return http.StatusInternalServerError, msg
}

func renderJSONError(w http.ResponseWriter, r *http.Request, err error) {
	code, msg := extractStatus(err)
	if code >= 500 {
		reportError(r.Context(), err, false)
		msg = "Internal Server Error"
	}

	renderJSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: msg,
	})
}
