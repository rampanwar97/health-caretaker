package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"health-caretaker/pkg/logger"
)

// Server represents an HTTP server with graceful shutdown
type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
	name       string
}

// New creates a new server instance
func New(addr string, handler http.Handler, name string, log *logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: log,
		name:   name,
	}
}

// Start starts the server and handles graceful shutdown
func (s *Server) Start() error {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting %s server on %s", s.name, s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start %s server: %v", s.name, err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	s.logger.Info("Shutting down %s server...", s.name)

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown %s server: %v", s.name, err)
	}

	s.logger.Info("%s server stopped", s.name)
	return nil
}

// StartBackground starts the server in the background without graceful shutdown
func (s *Server) StartBackground() error {
	s.logger.Info("Starting %s server on %s", s.name, s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start %s server: %v", s.name, err)
		}
	}()
	return nil
}

// Stop stops the server gracefully
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping %s server...", s.name)
	return s.httpServer.Shutdown(ctx)
}
