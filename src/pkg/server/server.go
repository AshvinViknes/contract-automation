package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"smart-contract-automation/src/pkg/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	router *chi.Mux
}

func NewServer(router *chi.Mux, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		},
		router: router,
	}
}

func (s *Server) Start() error {
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Server is starting", zap.String("port", s.server.Addr))
		serverErrors <- s.server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info("Server is shutting down", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
