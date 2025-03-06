package main

import (
	"log"
	"net/http"
	"time"
)

// CustomResponseWriter wraps http.ResponseWriter to capture status code and size.
type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (cw *CustomResponseWriter) WriteHeader(code int) {
	cw.statusCode = code
	cw.ResponseWriter.WriteHeader(code)
}

func (cw *CustomResponseWriter) Write(b []byte) (int, error) {
	size, err := cw.ResponseWriter.Write(b)
	cw.size += size
	return size, err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter
		cw := &CustomResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default to 200 unless changed
		}

		start := time.Now()
		log.Printf("Received request: %s %s", r.Method, r.RequestURI)

		next.ServeHTTP(cw, r) // pass our custom writer instead of "w"

		duration := time.Since(start)
		log.Printf("Response for %s %s => %d (%d bytes) took %v",
			r.Method, r.RequestURI, cw.statusCode, cw.size, duration)
	})
}
