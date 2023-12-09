package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type LoggerMiddleware struct {
	logger *log.Logger
}

func NewLoggerMiddleware(logger *log.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (l *LoggerMiddleware) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(rw, r)

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		l.logger.Printf(
			"\"%s %s://%s%s %s\" from %s - %d %dB in %v",
			r.Method,
			scheme,
			r.Host,
			r.RequestURI,
			r.Proto,
			r.RemoteAddr,
			rw.Status(),
			rw.BytesWritten(),
			time.Since(startTime),
		)
	})
}

func (l *LoggerMiddleware) Apply(next http.Handler) http.Handler {
	return l.logRequest(next)
}
