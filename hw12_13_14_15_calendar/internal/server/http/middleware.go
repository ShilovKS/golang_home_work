package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var globalLogger Logger

func SetLogger(logger Logger) {
	globalLogger = logger
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &StatusRecorder{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		clientIP := strings.Split(r.RemoteAddr, ":")[0]
		logMessage := fmt.Sprintf(
			"%s [%s] %s %s %s %d %d \"%s\"",
			clientIP,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			rw.statusCode,
			time.Since(start).Milliseconds(),
			r.UserAgent(),
		)
		globalLogger.Info(logMessage)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *StatusRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
