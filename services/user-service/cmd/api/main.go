package main

import (
	"log"

	"github.com/morshedulmunna/go-microservice/pkg/config"
	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/server"
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

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		logger.Fatal("Failed to create server", zap.Error(err))
	}

	// Start server
	if err := srv.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	// Wait for termination signal
	srv.WaitForSignal()

	// Stop server
	if err := srv.Stop(); err != nil {
		logger.Fatal("Failed to stop server", zap.Error(err))
	}
}
