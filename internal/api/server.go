package api

import (
	"context"
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"net/http"
	"time"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/api/routes/handlers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type DependencyConfig struct {
	port                string
	userController      interfaces.UserController
	appController       interfaces.AppController
	dockerController    interfaces.DockerController
	authController      interfaces.AuthController
	serverController    interfaces.ServerController
	webSocketController interfaces.WsController
	routeController     interfaces.RouteController
	jwt                 *middleware.JWT
	rateLimiter         *middleware.RateLimiter
}

func NewDependencyConfig(port string, userController interfaces.UserController,
	appController interfaces.AppController, dockerController interfaces.DockerController,
	authController interfaces.AuthController, jwt *middleware.JWT, serverController interfaces.ServerController,
	wsController interfaces.WsController, rateLimiter *middleware.RateLimiter, routeController interfaces.RouteController) *DependencyConfig {
	return &DependencyConfig{
		port:                port,
		userController:      userController,
		appController:       appController,
		dockerController:    dockerController,
		authController:      authController,
		serverController:    serverController,
		webSocketController: wsController,
		routeController:     routeController,
		jwt:                 jwt,
		rateLimiter:         rateLimiter,
	}
}

type Server struct {
	config *DependencyConfig
	server *http.Server
	router *routes.Router
}

func NewServer(cfg *DependencyConfig) *Server {
	return &Server{
		config: cfg,
		router: routes.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.SetupMiddleware()
	s.SetupRoutes()
	s.server = &http.Server{
		Addr:         ":" + s.config.port,
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s.server.ListenAndServe()
}

func (s *Server) SetupRoutes() {
	authHandler := handlers.NewAuthHandler(s.config.userController, s.config.authController, s.config.jwt, s.config.rateLimiter)
	userHandler := handlers.NewUserHandler(s.config.userController, s.config.jwt)
	appHandler := handlers.NewAppAppHandler(s.config.appController, s.config.dockerController, s.config.jwt)
	serverHandler := handlers.NewServerHandlers(s.config.serverController, s.config.jwt)
	wsHandler := handlers.NewWebsocketHandler(s.config.webSocketController, s.config.jwt)
	routeHandler := handlers.NewRouteHandlers(s.config.routeController)
	authHandler.SetupAuthHandlers(*s.router)
	appHandler.SetupAppHandlers(*s.router)
	wsHandler.SetupWebsocketHandlers(*s.router)
	serverHandler.SetupServerHandlers(*s.router)
	userHandler.SetupUserHandlers(*s.router)
	routeHandler.SetupRouteHandler(*s.router)

}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) SetupMiddleware() {
	s.router.USE(middleware.Logger)
	s.router.USE(middleware.CorsHandler)
	s.router.USE(middleware.ErrorHandler)
}
