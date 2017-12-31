package app

import (
	"log"
	"net/http"
	"time"
)

// Logger is logging middleware for the API
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s - %s (%v)\n", r.Method, r.URL.Path, time.Since(startTime))
	})
}
