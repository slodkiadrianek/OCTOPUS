package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services"

	"github.com/slodkiadrianek/octopus/internal/api"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

func main() {
	loggerService := logger.NewLogger("./logs", "02.01.2006")
	cfg := config.SetConfig("./.env")
	db := config.NewDb(cfg.DbLink)
	userRepository := repository.NewUserRepository(db.DbConnection, loggerService)
	userService := services.NewUserService(loggerService, userRepository)
	userController := controllers.NewUserController(userService)
	dependenciesConfig := api.NewDependencyConfig(cfg.Port, userController)
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
