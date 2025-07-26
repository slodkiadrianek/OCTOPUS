package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type Config struct {
	Port   string
	Logger *logger.Logger
}

type Server struct {
	config *config.Env
	server *http.Server
	router *routes.Router
}

func NewServer(cfg *config.Env) *Server {
	return &Server{
		config: cfg,
		router: routes.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.SetupMiddlware()
	s.SetupRoutes()
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
	s.router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hi")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hi from server"))
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) SetupMiddlware() {
	s.router.Use(middleware.ErrorHandler)
	s.router.Use(middleware.CorsHandler)
}
