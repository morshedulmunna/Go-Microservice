package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/morshedulmunna/go-microservice/pkg/config"
	grpcServer "github.com/morshedulmunna/go-microservice/pkg/grpc"
	"github.com/morshedulmunna/go-microservice/pkg/health"
	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"github.com/morshedulmunna/go-microservice/pkg/middleware"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/handlers"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/repository"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/service"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	cfg     *config.Config
	logger  *zap.Logger
	handler *handlers.UserHandler
	server  *http.Server
}

// NewServer creates a new Server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize logger
	logger := logger.GetLogger()

	// Initialize repository
	repo, err := repository.NewUserRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %v", err)
	}

	// Initialize service
	svc := service.NewUserService(repo)

	// Initialize handler
	handler := handlers.NewUserHandler(svc, logger)

	// Create server
	return &Server{
		cfg:     cfg,
		logger:  logger,
		handler: handler,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	// Create router
	mux := http.NewServeMux()

	// Create health check handler
	healthHandler := health.NewHandler(s.cfg.App.Name, s.cfg.App.Version)

	// Register routes
	mux.HandleFunc("/api/v1/users", s.handler.HandleUsers())
	mux.HandleFunc("/api/v1/users/", s.handler.HandleUser())

	// Create final handler with middleware
	mainHandler := middleware.Chain(
		mux,
		middleware.Logger(s.logger),
		middleware.Recover(),
		middleware.CORS(s.cfg.CORS.AllowedOrigins),
		middleware.RateLimit(s.cfg.RateLimit.Requests, time.Duration(s.cfg.RateLimit.Duration)*time.Second),
	)

	// Create a new mux for combining health check and main handler
	finalMux := http.NewServeMux()
	finalMux.Handle("/health", healthHandler)
	finalMux.Handle("/", mainHandler)

	// Create server
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.HTTP.Port),
		Handler:      finalMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server
	go func() {
		s.logger.Info("Starting server",
			zap.String("app_name", s.cfg.App.Name),
			zap.String("version", s.cfg.App.Version),
			zap.Int("port", s.cfg.HTTP.Port),
		)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the server
func (s *Server) Stop() error {
	s.logger.Info("Stopping server")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %v", err)
	}

	return nil
}

// WaitForSignal waits for termination signal
func (s *Server) WaitForSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	s.logger.Info("Received signal", zap.String("signal", sig.String()))
}

// RegisterGRPCServices registers gRPC services
func (s *Server) RegisterGRPCServices(srv *grpcServer.Server) {
	// TODO: Register gRPC services
}

// Close closes server resources
func (s *Server) Close() error {
	// TODO: Close any resources (e.g., database connections)
	return nil
}
