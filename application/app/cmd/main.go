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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// --- метрики ---
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// --- middleware ---
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		// пропускаем /metrics чтобы не засорять статистику
		if r.URL.Path != "/metrics" {
			httpRequestsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				fmt.Sprintf("%d", rw.statusCode),
			).Inc()
			httpRequestDuration.WithLabelValues(
				r.Method,
				r.URL.Path,
			).Observe(duration)
		}

		log.Printf("method=%s path=%s status=%d duration=%.3fs",
			r.Method, r.URL.Path, rw.statusCode, duration)
	})
}

// --- handlers ---
func databaseDSN() string {
	host := os.Getenv("DB_HOST")
	port := envOrDefault("DB_PORT", "5432")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslMode := envOrDefault("DB_SSLMODE", "require")
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, name, user, password, sslMode,
	)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
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