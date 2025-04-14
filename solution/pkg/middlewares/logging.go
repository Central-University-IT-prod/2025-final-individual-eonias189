package middlewares

import (
	"advertising/pkg/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggerProvider(l *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithCtx(r.Context(), l)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.FromCtx(r.Context())
			lrw := newLoggingResponseWriter(w)
			defer func(start time.Time) {
				l.Info(
					"request completed",
					zap.String("method", r.Method),
					zap.String("uri", r.RequestURI),
					zap.Int("status_code", lrw.statusCode),
					zap.Duration("elapsed_ms", time.Since(start)*1000),
				)
			}(time.Now())

			next.ServeHTTP(lrw, r)
		})
	}

}
