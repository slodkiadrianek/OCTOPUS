package main

import (
	"context"
	"fmt"
	"time"

	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
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
	appRepository := repository.NewAppRepository(db.DbConnection, loggerService)
	routeRepository := repository.NewRouteRepository(db.DbConnection, loggerService)
	appService := services.NewAppService(appRepository, loggerService, cacheService, cfg.DockerHost, routeRepository)
	serverService := services.NewServerService(loggerService, cacheService)
	ctx := context.Background()
	ticker(ctx, appService, serverService, loggerService)
}

func ticker(ctx context.Context, appService *services.AppService, serverService *services.ServerService, logger *utils.Logger) {
	period := 1 * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			appsToSendNotification, err := appService.CheckAppsStatus(ctx)
			fmt.Println()
			if err != nil {
				logger.Error("Something went wrong during checking statuses of apps", err)
			}
			err = appService.SendNotifications(ctx, appsToSendNotification)
			if err != nil {
				logger.Error("Something went wrong during checking statuses of apps", err)
			}
			logger.Info("Successfully inserted data about apps statuses")
			err = serverService.InsertServerMetrics(ctx)
			if err != nil {
				logger.Warn("Something went wrong during inserting data about server metrics", err)
			}
			logger.Info("Successfully inserted data about server status")
			err = appService.CheckRoutesStatus(ctx)
			if err != nil {
				logger.Warn("Something went wrong during checking statuses of the routes", err)
			}
		case <-ctx.Done():
			logger.Info("Status checked stopped")
			return
		}
	}
}
