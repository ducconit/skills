// Package examples demonstrates the Functional Options pattern in Go.
package examples

import (
	"log"
	"time"
)

// ─── Functional Options Pattern ─────────────────────────────────
// Cho phép tạo struct với config linh hoạt, giữ API sạch,
// dễ extend mà không breaking changes.

// Server represents an HTTP server with configurable options.
type Server struct {
	host         string
	port         int
	timeout      time.Duration
	maxConns     int
	tls          bool
	logger       *log.Logger
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// Option defines a functional option for Server.
type Option func(*Server)

// WithHost sets the server host.
func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// WithPort sets the server port.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithTimeout sets the general timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// WithMaxConns sets the maximum number of connections.
func WithMaxConns(maxConns int) Option {
	return func(s *Server) {
		s.maxConns = maxConns
	}
}

// WithTLS enables TLS.
func WithTLS(enabled bool) Option {
	return func(s *Server) {
		s.tls = enabled
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithReadWriteTimeouts sets both read and write timeouts.
func WithReadWriteTimeouts(read, write time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = read
		s.writeTimeout = write
	}
}

// NewServer creates a Server with sensible defaults and optional overrides.
func NewServer(opts ...Option) *Server {
	// Sensible defaults
	s := &Server{
		host:         "0.0.0.0",
		port:         8080,
		timeout:      30 * time.Second,
		maxConns:     100,
		tls:          false,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// ─── Usage Examples ─────────────────────────────────────────────

func ExampleFunctionalOptions() {
	// Default server
	s1 := NewServer()
	_ = s1

	// Custom server
	s2 := NewServer(
		WithPort(9090),
		WithTLS(true),
		WithTimeout(60*time.Second),
		WithMaxConns(500),
	)
	_ = s2

	// Production server
	s3 := NewServer(
		WithHost("0.0.0.0"),
		WithPort(443),
		WithTLS(true),
		WithMaxConns(10000),
		WithReadWriteTimeouts(10*time.Second, 30*time.Second),
	)
	_ = s3
}

// ─── Functional Options with Validation ─────────────────────────
// Variant: options return error để validate config.

type ValidatedOption func(*Server) error

func WithValidatedPort(port int) ValidatedOption {
	return func(s *Server) error {
		if port < 1 || port > 65535 {
			return &InvalidOptionError{
				Option: "port",
				Value:  port,
				Reason: "must be between 1 and 65535",
			}
		}
		s.port = port
		return nil
	}
}

type InvalidOptionError struct {
	Option string
	Value  interface{}
	Reason string
}

func (e *InvalidOptionError) Error() string {
	return "invalid option " + e.Option + ": " + e.Reason
}

func NewServerValidated(opts ...ValidatedOption) (*Server, error) {
	s := &Server{
		host:    "0.0.0.0",
		port:    8080,
		timeout: 30 * time.Second,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}
