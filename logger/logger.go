package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var DefaultLogger func(next http.Handler) http.Handler

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	logger := new(Logger)
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Logger = l
	DefaultLogger = logger.RequestLogger()

	return logger
}

func (l *Logger) LoggerMiddleware(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

func (l *Logger) RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Logger.Info("Request Middleware", "Status", ww.Status(), "Bytes", ww.BytesWritten(), "Method", r.Method, "ResponseTime", time.Since(t1).Milliseconds())
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(f)
	}
}
