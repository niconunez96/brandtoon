package telemetry

import (
	sharedhttp "brandtoonapi/bounded_contexts/shared/infra/http"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

var telemetryOnce sync.Once
var telemetry *Telemetry

type Telemetry struct {
	logger Logger
}

func NewTelemetry() *Telemetry {
	telemetryOnce.Do(func() {
		telemetry = &Telemetry{
			logger: newLogger(),
		}
	})
	return telemetry
}

func (t *Telemetry) HttpLoggerMiddleware() sharedhttp.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				t2 := time.Now()
				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					t.logger.LogError(
						"log system error",
						errors.New("internal error"),
						map[string]any{
							"recover_info": rec,
							"debug_stack":  debug.Stack(),
						},
					)
					http.Error(
						ww,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}

				// log end request
				t.logger.LogInfo(
					fmt.Sprintf("incoming_request %s", r.URL.Path),
					map[string]any{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					},
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

type attr struct {
	key   string
	value any
}

func Attr(key string, value any) attr {
	return attr{key, value}
}

func buildAttrMap(attrs []attr) map[string]any {
	attrsMap := make(map[string]any)
	for _, attr := range attrs {
		attrsMap[attr.key] = attr.value
	}
	return attrsMap
}

func LogDebug(msg string, attrs ...attr) {
	attrsMap := buildAttrMap(attrs)
	NewTelemetry().logger.LogInfo(msg, attrsMap)
}

func LogInfo(msg string, attrs ...attr) {
	attrsMap := buildAttrMap(attrs)
	NewTelemetry().logger.LogInfo(msg, attrsMap)
}

func LogError(msg string, err error, attrs ...attr) {
	attrsMap := buildAttrMap(attrs)
	NewTelemetry().logger.LogError(msg, err, attrsMap)
}
