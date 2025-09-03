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
	Port string
	UserController *controllers.UserController
	AuthController *controllers.AuthController
	JWT            *middleware.JWT
}

func NewDependencyConfig(port string, userController *controllers.UserController, authController *controllers.AuthController, jwt *middleware.JWT) *DependencyConfig {
	return &DependencyConfig{
		Port:           port,
		UserController: userController,
		AuthController: authController,
		JWT:            jwt,
		// CacheService:   cacheService,
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
	authHandler := handlers.NewAuthHandler(s.Config.UserController, s.Config.AuthController, s.Config.JWT)
	userHandler := handlers.NewUserHandler(s.Config.UserController, s.Config.JWT)
	authHandler.SetupAuthHandlers(*s.Router)
	userHandler.SetupUserHandlers(*s.Router)
	//usersApi := s.router.Group("/users")
	//usersApi.GET("/us", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("Hi")
	//	w.WriteHeader(http.StatusOK)
	//	w.Write([]byte("Hi from server"))
	//})
	//s.router.GET("/users/:userIdd", func(w http.ResponseWriter, r *http.Request) {
	//	userId, err := utils.ReadBody[map[string]string](r)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(userId)
	//	w.WriteHeader(http.StatusOK)
	//	w.Write([]byte(`Body readed successfully`))
	//})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) SetupMiddleware() {
	s.Router.USE(middleware.Logger)
	s.Router.USE(middleware.CorsHandler)
	s.Router.USE(middleware.ErrorHandler)
}
