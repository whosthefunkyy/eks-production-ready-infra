package main

import (
	"fmt"
	"log"
	"net/http"	
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "From Api, hii")
	})
	mux.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 500000000; i++ {
		_ = i * i
	}

	fmt.Fprintln(w, "cpu done")
})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.Handle("/metrics", promhttp.Handler())

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggingMiddleware(mux)))
}