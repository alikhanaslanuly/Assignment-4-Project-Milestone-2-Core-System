package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// AuthMockMiddleware simulates an authentication check
func AuthMockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple mock: check if "Authorization" header exists
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Auth failed: Missing Authorization header for %s", r.URL.Path)
			// For some public routes (like /register or /health), we might want to skip this.
			// But for a mock, let's just log it or apply it to specific routes in main.go
		}
		next.ServeHTTP(w, r)
	})
}
