package api

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/controllers"

	// "fmt"
	"net/http"
	"time"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/api/routes/handlers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type DependencyConfig struct {
	Port                string
	UserController      *controllers.UserController
	AppController       *controllers.AppController
	DockerController    *controllers.DockerController
	AuthController      *controllers.AuthController
	ServerController    *controllers.ServerController
	WebSocketController *controllers.WsController
	JWT                 *middleware.JWT
	RateLimiter         *middleware.RateLimiter
}

func NewDependencyConfig(port string, userController *controllers.UserController,
	appController *controllers.AppController, dockerController *controllers.DockerController,
	authController *controllers.AuthController, jwt *middleware.JWT, serverController *controllers.ServerController,
	wsController *controllers.WsController, rateLimiter *middleware.RateLimiter) *DependencyConfig {
	return &DependencyConfig{
		Port:                port,
		UserController:      userController,
		AppController:       appController,
		DockerController:    dockerController,
		AuthController:      authController,
		ServerController:    serverController,
		WebSocketController: wsController,
		JWT:                 jwt,
		RateLimiter:         rateLimiter,
	}
}

type Server struct {
	Config *DependencyConfig
	server *http.Server
	Router *routes.Router
}

func NewServer(cfg *DependencyConfig) *Server {
	return &Server{
		Config: cfg,
		Router: routes.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.SetupMiddleware()
	s.SetupRoutes()
	s.server = &http.Server{
		Addr:         ":" + s.Config.Port,
		Handler:      s.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s.server.ListenAndServe()
}

func (s *Server) SetupRoutes() {
	authHandler := handlers.NewAuthHandler(s.Config.UserController, s.Config.AuthController, s.Config.JWT, s.Config.RateLimiter)
	userHandler := handlers.NewUserHandler(s.Config.UserController, s.Config.JWT)
	appHandler := handlers.NewAppAppHandler(s.Config.AppController, s.Config.DockerController, s.Config.JWT)
	serverHandler := handlers.NewServerHandlers(s.Config.ServerController, s.Config.JWT)
	wsHandler := handlers.NewWebsocketHandler(s.Config.WebSocketController, s.Config.JWT)
	authHandler.SetupAuthHandlers(*s.Router)
	appHandler.SetupAppHandlers(*s.Router)
	wsHandler.SetupWebsocketHandlers(*s.Router)
	serverHandler.SetupServerHandlers(*s.Router)
	userHandler.SetupUserHandlers(*s.Router)

}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) SetupMiddleware() {
	s.Router.USE(middleware.Logger)
	s.Router.USE(middleware.CorsHandler)
	s.Router.USE(middleware.ErrorHandler)
}
