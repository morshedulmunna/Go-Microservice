package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/morshedulmunna/go-microservice/pkg/config"
	"github.com/morshedulmunna/go-microservice/pkg/health"
	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"github.com/morshedulmunna/go-microservice/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.GetLogger()
	defer logger.Sync()

	// Create router
	mux := http.NewServeMux()

	// Create health check handler
	healthHandler := health.NewHandler(cfg.App.Name, cfg.App.Version)

	// Add middleware
	mainHandler := middleware.Chain(
		mux,
		middleware.Logger(logger),
		middleware.Recover(),
		middleware.CORS(cfg.CORS.AllowedOrigins),
		middleware.RateLimit(cfg.RateLimit.Requests, time.Duration(cfg.RateLimit.Duration)*time.Second),
	)

	// Create a new mux for combining health check and main handler
	finalMux := http.NewServeMux()
	finalMux.Handle("/health", healthHandler)
	finalMux.Handle("/", mainHandler)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      finalMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server
	go func() {
		logger.Info("Starting order service...",
			zap.String("app_name", cfg.App.Name),
			zap.String("version", cfg.App.Version),
			zap.Int("port", cfg.HTTP.Port),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	logger.Info("Shutting down server...")
	if err := server.Close(); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}
}
