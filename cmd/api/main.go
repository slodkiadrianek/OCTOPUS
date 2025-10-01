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
	loggerService.CreateLogger()
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
	userRepository := repository.NewUserRepository(db.DbConnection, loggerService)
	userService := services.NewUserService(loggerService, userRepository)
	appRepository := repository.NewAppRepository(db.DbConnection, loggerService)
	dockerRepository := repository.NewDockerRepository(db.DbConnection, loggerService)
	appService := services.NewAppService(appRepository, loggerService, cacheService, cfg.DockerHost)
	wsService := services.NewWsService(loggerService, cfg.DockerHost)
	serverService := services.NewServerService(loggerService, cacheService)
	dockerService := services.NewDockerService(dockerRepository, appRepository, loggerService, cfg.DockerHost)
	jwt := middleware.NewJWT(cfg.JWTSecret, loggerService, cacheService)
	authService := services.NewAuthService(loggerService, userRepository, jwt)
	authController := controllers.NewAuthController(authService, loggerService)
	userController := controllers.NewUserController(userService, loggerService)
	appController := controllers.NewAppController(appService, loggerService)
	dockerController := controllers.NewDockerController(dockerService, loggerService)
	serverController := controllers.NewServerController(loggerService, serverService)
	WsController := controllers.NewWsController(wsService, loggerService)

	dependenciesConfig := api.NewDependencyConfig(cfg.Port, userController, appController, dockerController,
		authController, jwt, serverController, WsController)
	server := api.NewServer(dependenciesConfig)

	go func() {
		loggerService.Info("Starting API server on port: " + cfg.Port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loggerService.Error("Failed to start server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		loggerService.Error("Server forced to shutdown:", err)
	}

	loggerService.Info("Server exited")
}
