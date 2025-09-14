package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type AppService struct {
	AppRepository *repository.AppRepository
	Logger        *logger.Logger
	CacheService  *config.CacheService
}

func NewAppService(appRepository *repository.AppRepository, logger *logger.Logger, cacheService *config.CacheService) *AppService {
	return &AppService{
		AppRepository: appRepository,
		Logger:        logger,
		CacheService:  cacheService,
	}
}

func (a *AppService) CreateApp(ctx context.Context, app schema.CreateApp, ownerId int) error {
	appDto := DTO.NewApp(app.Name, app.Description, app.DbLink, app.ApiLink, ownerId, app.DiscordWebhook, app.SlackWebhook)
	err := a.AppRepository.InsertApp(ctx, *appDto)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetApp(ctx context.Context, id int) (*models.App, error) {
	app, err := a.AppRepository.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) UpdateApp(ctx context.Context, id int, app schema.UpdateApp) error {
	appDto := DTO.NewUpdateApp(id, app.Name, app.Description, app.DbLink, app.ApiLink, app.DiscordWebhook, app.SlackWebhook)
	err := a.AppRepository.UpdateApp(ctx, *appDto)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) DeleteApp(ctx context.Context, id int) error {
	err := a.AppRepository.DeleteApp(ctx, id)
	if err != nil {
		return err
	}
	return nil
}