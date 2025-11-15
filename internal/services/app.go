package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type appRepository interface {
	InsertApp(ctx context.Context, app []DTO.App) error
	GetApp(ctx context.Context, id string, ownerId int) (*models.App, error)
	GetApps(ctx context.Context, ownerId int) ([]models.App, error)
	DeleteApp(ctx context.Context, id string, ownerId int) error
	GetAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error)
	GetAppsToCheck(ctx context.Context) ([]*models.AppToCheck, error)
	UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error
	InsertAppStatuses(ctx context.Context, appsStatuses []DTO.AppStatus) error
	GetUsersToSendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) ([]models.NotificationInfo, error)
}

type AppService struct {
	AppRepository           appRepository
	LoggerService           *utils.Logger
	appStatusService        *AppStatusService
	appNotificationsService *AppNotificationsService
	routeStatusService      *RouteStatusService
}

func NewAppService(appRepository appRepository, loggerService *utils.Logger,
	appStatusService *AppStatusService, appNotificationsService *AppNotificationsService,
	routeStatusService *RouteStatusService,
) *AppService {
	return &AppService{
		AppRepository:           appRepository,
		LoggerService:           loggerService,
		appStatusService:        appStatusService,
		appNotificationsService: appNotificationsService,
		routeStatusService:      routeStatusService,
	}
}

func (a *AppService) CreateApp(ctx context.Context, app DTO.CreateApp, ownerId int) error {
	GeneratedId, err := utils.GenerateID()
	if err != nil {
		return err
	}
	appDto := DTO.NewApp(GeneratedId, app.Name, app.Description, false, ownerId, app.IpAddress, app.Port)
	err = a.AppRepository.InsertApp(ctx, []DTO.App{*appDto})
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetApp(ctx context.Context, id string, ownerId int) (*models.App, error) {
	app, err := a.AppRepository.GetApp(ctx, id, ownerId)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) GetApps(ctx context.Context, ownerId int) ([]models.App, error) {
	apps, err := a.AppRepository.GetApps(ctx, ownerId)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (a *AppService) DeleteApp(ctx context.Context, appId string, ownerId int) error {
	err := a.AppRepository.DeleteApp(ctx, appId, ownerId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error {
	err := a.AppRepository.UpdateApp(ctx, appId, app, ownerId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetAppStatus(ctx context.Context, appId string, ownerId int) (DTO.AppStatus, error) {
	return a.appStatusService.getAppStatus(ctx, appId, ownerId)
}

func (a *AppService) CheckAppsStatus(ctx context.Context) ([]DTO.AppStatus, error) {
	return a.appStatusService.checkAppsStatus(ctx)
}

func (a *AppService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	return a.appNotificationsService.sendNotifications(ctx, appsStatuses)
}

func (a *AppService) CheckRoutesStatus(ctx context.Context) error {
	return a.routeStatusService.checkRoutesStatus(ctx)
}
