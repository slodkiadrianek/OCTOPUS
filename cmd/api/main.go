package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/slodkiadrianek/octopus/internal/api"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

func main() {
	loggerService := logger.NewLogger("./logs", "02.01.2006")
	cfg := config.SetConfig()
	server := api.NewServer(cfg)

	go func() {
		loggerService.Info("Starting API server on port" + cfg.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			loggerService.Error("Failed to start server", err)
		}
	}()

	loggerService.Info("Shutting down server...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// 10. Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		loggerService.Error("Server forced to shutdown:", err)
	}

	loggerService.Info("Server exited")
}
