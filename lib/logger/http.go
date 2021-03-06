package logger

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type HTTPLogger interface {
	ReadBody() error
	Add(f func(e *zerolog.Event))
	Write(t time.Time)
}

type httpLogger struct {
	l    *zerolog.Logger
	r    *http.Request
	ww   chi_middleware.WrapResponseWriter
	body []byte
	add  []func(e *zerolog.Event)
}

func NewHTTPLogger(r *http.Request, ww chi_middleware.WrapResponseWriter) HTTPLogger {
	return &httpLogger{
		l:    NewBaseLogger(),
		r:    r,
		ww:   ww,
		body: []byte{},
		add:  []func(e *zerolog.Event){},
	}
}

func (hl *httpLogger) ReadBody() error {
	if hl.r.Body != nil {
		body, err := ioutil.ReadAll(hl.r.Body)
		if err != nil {
			hl.add = append(hl.add, func(e *zerolog.Event) {
				e.Err(errors.Wrap(err, "read http body error"))
			})
			return err
		}
		hl.body = append(hl.body, body...)
		hl.r.Body = ioutil.NopCloser(
			bytes.NewBuffer(body),
		)
	}
	return nil
}

func (hl *httpLogger) Add(f func(e *zerolog.Event)) {
	hl.add = append(hl.add, f)
}

func (hl *httpLogger) Write(t time.Time) {
	e := hl.l.Log().Timestamp().
		Str("requestId", chi_middleware.GetReqID(hl.r.Context())).
		Dur("elapsed(ms)", time.Since(t)).
		Str("protocol", hl.r.Proto).
		Str("method", hl.r.Method).
		Str("uri", hl.r.URL.RequestURI()).
		Bytes("body", hl.body).
		Str("origin", hl.r.Header.Get("Origin")).
		Str("referer", hl.r.Referer()).
		Int("status", hl.ww.Status()).
		Str("agent", hl.r.UserAgent())

	for _, f := range hl.add {
		f(e)
	}

	if ip := clientIP(hl.r.Header); ip != "" {
		e.Str("client-ip", ip)
	} else if ip, _, err := net.SplitHostPort(strings.TrimSpace(hl.r.RemoteAddr)); err == nil {
		e.Str("client-ip", ip)
	}

	e.Str("cookies", fmt.Sprint(hl.r.Cookies()))
	e.Send()
}

var (
	trueClientIP          = http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor         = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP               = http.CanonicalHeaderKey("X-Real-IP")
	xEnvoyExternalAddress = http.CanonicalHeaderKey("X-Envoy-External-Address")
)

// clientIP returns the IP of the client.
// If a header identifying the real IP exists, the value of the header will be used.
func clientIP(h http.Header) string {
	if tcip := h.Get(trueClientIP); tcip != "" {
		return tcip
	} else if xrip := h.Get(xRealIP); xrip != "" {
		return xrip
	} else if xff := h.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		return xff[:i]
	} else if xeea := h.Get(xEnvoyExternalAddress); xeea != "" {
		return xeea
	}

	return ""
}

const (
	Key_HTTPLoggerCtx = "Key_HTTPLoggerCtx"
)

func HTTPLoggerCtx(ctx context.Context) HTTPLogger {
	v := ctx.Value(Key_HTTPLoggerCtx)
	if v, ok := v.(HTTPLogger); ok {
		return v
	}
	return nil
}

func SetHTTPLoggerCtx(ctx context.Context, hl HTTPLogger) context.Context {
	return context.WithValue(ctx, Key_HTTPLoggerCtx, hl)
}
