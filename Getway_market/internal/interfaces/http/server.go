package http

import (
	"net/http"
	"time"

	"Getway_market/internal/infrastructure/observability"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewServer(handler http.Handler, allowedOrigins []string) *http.Server {
	r := chi.NewRouter()
	r.Use(observability.Recover)
	r.Use(observability.RequestID)
	r.Use(observability.BodyLimit(10 << 20)) // 10MB
	r.Use(observability.AccessLog)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID", "Idempotency-Key"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// всё остальное — в прокси
	r.NotFound(handler.ServeHTTP)

	return &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
