package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"

	"github.com/slodkiadrianek/octopus/internal/api"
	"github.com/slodkiadrianek/octopus/internal/config"
)

func main() {
	loggerService := utils.NewLogger("./logs", "2006-01-02 15:04:05")
	loggerService.InitializeLogger()
	defer loggerService.Close()
	cfg, err := config.SetConfig("./.env")
	if err != nil {
		loggerService.Error("Failed to load config", err)
		return
	}
	err = cfg.Validate()
	if err != nil {
		loggerService.Error("Configuration validation failed", err)
		return
	}
	cacheService, err := config.NewCacheService(cfg.CacheLink)
	if err != nil {
		loggerService.Error("Failed to connect to cache", err)
		return
	}
	db, err := config.NewDb(cfg.DbLink, "postgres")
	if err != nil {
		loggerService.Error("Failed to connect to database", err)
		return
	}
	rateLimiter := middleware.NewRateLimiter(5, 1*time.Minute, 5*time.Minute, 10*time.Minute, loggerService)
	jwt := middleware.NewJWT(cfg.JWTSecret, loggerService, cacheService)
	// User
	userRepository := repository.NewUserRepository(db.DbConnection, loggerService)
	userService := services.NewUserService(loggerService, userRepository)
	userController := controllers.NewUserController(userService, loggerService)
	// Route
	routeRepository := repository.NewRouteRepository(db.DbConnection, loggerService)
	routeService := services.NewRouteService(loggerService, routeRepository)
	routeController := controllers.NewRouteController(routeService, loggerService)
	// App
	appRepository := repository.NewAppRepository(db.DbConnection, loggerService)
	appService := services.NewAppService(appRepository, loggerService, cacheService, cfg.DockerHost, routeRepository)
	appController := controllers.NewAppController(appService, loggerService)
	// webSocket
	wsService := services.NewWsService(loggerService, cfg.DockerHost)
	wsController := controllers.NewWsController(wsService, loggerService)
	// Server
	serverService := services.NewServerService(loggerService, cacheService)
	serverController := controllers.NewServerController(loggerService, serverService)
	// Docker
	dockerService := services.NewDockerService(appRepository, loggerService, cfg.DockerHost)
	dockerController := controllers.NewDockerController(dockerService, loggerService)
	// Auth
	authService := services.NewAuthService(loggerService, userRepository, jwt)
	authController := controllers.NewAuthController(authService, loggerService)

	dependenciesConfig := api.NewDependencyConfig(cfg.Port, userController, appController, dockerController,
		authController, jwt, serverController, wsController, rateLimiter, routeController)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	server := api.NewServer(dependenciesConfig)

	go rateLimiter.CleanWorker(ctx)
	go func() {
		loggerService.Info("Starting API server on port: " + cfg.Port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loggerService.Error("Failed to start server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		loggerService.Error("Server forced to shutdown:", err)
	}

	loggerService.Info("Server exited")
}
