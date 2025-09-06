package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
<<<<<<< HEAD
	"github.com/slodkiadrianek/octopus/internal/models"
=======
>>>>>>> a4f4bb342f74a1e297be363d81262025c784bffa
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
<<<<<<< HEAD
	appDto := DTO.NewApp(app.Name, app.Description, app.DbLink, app.ApiLink, ownerId, app.DiscordWebhook, app.SlackWebhook)
=======
	appDto := DTO.NewApp(app.Name, app.DbLink, app.ApiLink, ownerId, app.DiscordWebhook, app.SlackWebhook)
>>>>>>> a4f4bb342f74a1e297be363d81262025c784bffa
	err := a.AppRepository.InsertApp(ctx, *appDto)
	if err != nil {
		return err
	}
	return nil
}

<<<<<<< HEAD
func (a *AppService) GetApp(ctx context.Context, id int) (*models.App, error) {
=======
func (a *AppService) GetApp(ctx context.Context, id int) (*DTO.App, error) {
>>>>>>> a4f4bb342f74a1e297be363d81262025c784bffa
	app, err := a.AppRepository.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) UpdateApp(ctx context.Context, id int, app schema.UpdateApp) error {
<<<<<<< HEAD
	appDto := DTO.NewUpdateApp(id, app.Name, app.Description, app.DbLink, app.ApiLink, app.DiscordWebhook, app.SlackWebhook)
	err := a.AppRepository.UpdateApp(ctx, *appDto)
=======
	appDto := DTO.NewApp(app.Name, app.DbLink, app.ApiLink, app.OwnerId, app.DiscordWebhook, app.SlackWebhook)
	err := a.AppRepository.UpdateApp(ctx, id, *appDto)
>>>>>>> a4f4bb342f74a1e297be363d81262025c784bffa
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
