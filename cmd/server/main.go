package main

import (
	"example-app/internal/debug"
	"example-app/internal/health"
	"example-app/internal/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	healthHandler := http.HandlerFunc(health.Handler)
	debugErrorHandler := http.HandlerFunc(debug.ErrorHandler)

	mux.Handle("GET /health", middleware.Logging(middleware.Metrics("/health", healthHandler)))
	mux.Handle("GET /debug/error", middleware.Logging(middleware.Metrics("/debug/error", debugErrorHandler)))
	mux.Handle("GET /metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("server starting", "port", port)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}

}
