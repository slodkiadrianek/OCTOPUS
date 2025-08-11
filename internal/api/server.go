package api

import (
	"context"
	"fmt"

	// "fmt"
	"net/http"
	"time"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/api/routes/handlers"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/utils"
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
	s.SetupMiddleware()
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
	handlers.SetupAuthHadnlers(s.router)
	usersApi := s.router.Group("/users")
	usersApi.GET("/us", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hi")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hi from server"))
	})
	s.router.GET("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.ReadParam(r, "userId")
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userId))
		w.Write([]byte(userId))
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) SetupMiddleware() {
	s.router.USE(middleware.Logger)
	s.router.USE(middleware.CorsHandler)
	s.router.USE(middleware.ErrorHandler)
}
