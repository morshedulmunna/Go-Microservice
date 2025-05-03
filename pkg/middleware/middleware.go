package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"go.uber.org/zap"
)

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Chain chains multiple middleware functions together
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		if m != nil {
			handler = m(handler)
		}
	}
	return handler
}

// Logger creates a middleware for logging HTTP requests
func Logger(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}

// Recover creates a middleware for recovering from panics
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS creates a middleware for handling CORS
func CORS(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if allowed == "*" || allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
					break
				}
			}
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit creates a middleware for rate limiting requests
func RateLimit(requests int, duration time.Duration) Middleware {
	// TODO: Implement rate limiting using a token bucket or similar algorithm
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout middleware adds a timeout to the request context
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)
			done := make(chan bool)

			go func() {
				next.ServeHTTP(w, r)
				done <- true
			}()

			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				logger.Error("Request timeout",
					zap.String("path", r.URL.Path),
					zap.Duration("timeout", timeout),
				)
				return
			case <-done:
				return
			}
		})
	}
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// AuthMiddleware is a placeholder for JWT authentication
// This should be implemented based on your authentication requirements
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: Implement JWT token validation
		// This is where you would validate the token and set user information in the context

		next.ServeHTTP(w, r)
	})
}
