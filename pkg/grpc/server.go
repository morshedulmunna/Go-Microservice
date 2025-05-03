package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server represents a gRPC server
type Server struct {
	server     *grpc.Server
	port       int
	logger     *zap.Logger
	opts       []grpc.ServerOption
	services   []Service
	healthSrv  *health.Server
	reflection bool
}

// Service represents a gRPC service
type Service interface {
	Register(*grpc.Server)
}

// ServerOption represents a server option
type ServerOption func(*Server)

// NewServer creates a new gRPC server
func NewServer(port int, options ...ServerOption) *Server {
	s := &Server{
		port:       port,
		logger:     logger.GetLogger(),
		opts:       make([]grpc.ServerOption, 0),
		services:   make([]Service, 0),
		healthSrv:  health.NewServer(),
		reflection: false,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

// WithTLS enables TLS for the server
func WithTLS(certFile, keyFile string) ServerOption {
	return func(s *Server) {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			s.logger.Fatal("Failed to load TLS credentials", zap.Error(err))
		}
		s.opts = append(s.opts, grpc.Creds(creds))
	}
}

// WithUnaryInterceptors adds unary interceptors
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.opts = append(s.opts, grpc.ChainUnaryInterceptor(interceptors...))
	}
}

// WithStreamInterceptors adds stream interceptors
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.opts = append(s.opts, grpc.ChainStreamInterceptor(interceptors...))
	}
}

// WithKeepalive adds keepalive parameters
func WithKeepalive(params keepalive.ServerParameters) ServerOption {
	return func(s *Server) {
		s.opts = append(s.opts, grpc.KeepaliveParams(params))
	}
}

// WithReflection enables server reflection
func WithReflection() ServerOption {
	return func(s *Server) {
		s.reflection = true
	}
}

// RegisterService registers a gRPC service
func (s *Server) RegisterService(service Service) {
	s.services = append(s.services, service)
}

// Start starts the gRPC server
func (s *Server) Start() error {
	// Add default keepalive parameters
	s.opts = append(s.opts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Minute,
		Time:                  5 * time.Minute,
		Timeout:               20 * time.Second,
	}))

	// Create gRPC server
	s.server = grpc.NewServer(s.opts...)

	// Register services
	for _, service := range s.services {
		service.Register(s.server)
	}

	// Register health service
	healthpb.RegisterHealthServer(s.server, s.healthSrv)

	// Enable reflection if requested
	if s.reflection {
		reflection.Register(s.server)
	}

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.logger.Info("Starting gRPC server",
		zap.Int("port", s.port),
		zap.Bool("reflection", s.reflection),
	)

	return s.server.Serve(lis)
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	if s.server != nil {
		s.logger.Info("Stopping gRPC server")
		s.server.GracefulStop()
	}
}

// UnaryServerInterceptor creates a new unary server interceptor
func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Log request
		logger.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		// Handle request
		resp, err := handler(ctx, req)

		// Log response
		logger.Info("gRPC response",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}

// StreamServerInterceptor creates a new stream server interceptor
func StreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Log stream start
		logger.Info("gRPC stream started",
			zap.String("method", info.FullMethod),
		)

		// Handle stream
		err := handler(srv, ss)

		// Log stream end
		logger.Info("gRPC stream ended",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return err
	}
}

// SetServingStatus sets the serving status of the health check
func (s *Server) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.healthSrv.SetServingStatus(service, status)
}

// GetHealthServer returns the health server
func (s *Server) GetHealthServer() *health.Server {
	return s.healthSrv
}
