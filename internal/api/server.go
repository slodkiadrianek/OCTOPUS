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
	Port           string
	//Logger         *logger.Logger
	UserController *controllers.UserController
}

func NewDependencyConfig(port string, userController *controllers.UserController) *DependencyConfig{
	return &DependencyConfig{
		Port: port,
		UserController: userController,
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
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s.server.ListenAndServe()
}

func (s *Server) SetupRoutes() {
	authHandler := handlers.NewAuthHandler(s.config.UserController)
	authHandler.SetupAuthHandlers(*s.router)
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
	s.router.USE(middleware.Logger)
	s.router.USE(middleware.CorsHandler)
	s.router.USE(middleware.ErrorHandler)
}
