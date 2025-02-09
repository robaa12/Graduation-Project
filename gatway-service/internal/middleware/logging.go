package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}
