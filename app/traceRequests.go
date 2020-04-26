package app

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func TranceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request -> method: %s, addr: %s, path: %s\n", r.Method, r.RemoteAddr, r.URL.Path)
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		log.Printf("response -> statusCode: %d, headers: %s\n", lrw.statusCode, lrw.Header())
	})
}
