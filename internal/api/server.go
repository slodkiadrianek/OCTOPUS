package api

import (
	"context"
	"net/http"
	"time"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type Config struct {
	Port   string
	Logger *logger.Logger
}

type Server struct {
	config *Config
	server *http.Server
	router *routes.Router
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s.server.ListenAndServe()
}

func (s *Server) SetupRoutes() {
	s.router = routes.NewRouter()
	s.router.Get("/users/:id")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) SetupMiddlware() {
	s.router.Use(middleware.ErrorHandler)
	s.router.Use(middleware.CorsHandler)
}
