package main

import (
	"context"
	"time"

	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

func main() {
	loggerService := logger.NewLogger("./logs", "02.01.2006")
	cfg, err := config.SetConfig("./.env")
	if err != nil {
		loggerService.Error("Failed to load config", err)
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
	appRepository := repository.NewAppRepository(db.DbConnection, loggerService)
	appService := services.NewAppService(appRepository, loggerService, cacheService, cfg.DockerHost)
	ctx := context.Background()
	runStatusChecker(ctx, appService, loggerService)
}

func runStatusChecker(ctx context.Context, appService *services.AppService, logger *logger.Logger) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			appsToSendNotification, err := appService.CheckAppsStatus(ctx)
			if err != nil {
				logger.Error("Something went wrong during checking statuses of apps", err)
			}
			err = appService.SendNotifications(ctx, appsToSendNotification)
			if err != nil {
				logger.Error("Something went wrong during checking statuses of apps", err)
			}
			logger.Info("Succesfully inserted data about apps statuses")
		case <-ctx.Done():
			logger.Info("Status checked stopped")
			return
		}
	}
}
