package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}

func databaseDSN() string {
	host := os.Getenv("DB_HOST")
	port := envOrDefault("DB_PORT", "5432")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host,
		port,
		name,
		user,
		password,
	)
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func dbHealthHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", databaseDSN())
	if err != nil {
		http.Error(w, "db open failed", http.StatusInternalServerError)
		log.Printf("db open failed: %v", err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		http.Error(w, "db ping failed", http.StatusServiceUnavailable)
		log.Printf("db ping failed: %v", err)
		return
	}

	fmt.Fprintln(w, "db ok")
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
	mux.HandleFunc("/db-health", dbHealthHandler)

	mux.Handle("/metrics", promhttp.Handler())

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggingMiddleware(mux)))
}
