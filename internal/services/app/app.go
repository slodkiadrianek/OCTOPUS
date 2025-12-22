package servicesApp

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AppService struct {
	appRepository           interfaces.AppRepository
	loggerService           utils.LoggerService
	appStatusService        interfaces.AppStatusService
	appNotificationsService interfaces.AppNotificationsService
	routeStatusService      interfaces.RouteStatusService
}

func NewAppService(appRepository interfaces.AppRepository, loggerService utils.LoggerService,
	appStatusService interfaces.AppStatusService, appNotificationsService interfaces.AppNotificationsService,
	routeStatusService interfaces.RouteStatusService,
) *AppService {
	return &AppService{
		appRepository:           appRepository,
		loggerService:           loggerService,
		appStatusService:        appStatusService,
		appNotificationsService: appNotificationsService,
		routeStatusService:      routeStatusService,
	}
}

func (a *AppService) CreateApp(ctx context.Context, app DTO.CreateApp, ownerID int) error {
	generatedID, err := utils.GenerateID()
	if err != nil {
		return err
	}

	appDto := DTO.NewApp(generatedID, app.Name, app.Description, false, ownerID, app.IPAddress, app.Port)
	err = a.appRepository.InsertApp(ctx, []DTO.App{*appDto})
	if err != nil {
		return err
	}

	return nil
}

func (a *AppService) GetApp(ctx context.Context, appID string, ownerID int) (*models.App, error) {
	app, err := a.appRepository.GetApp(ctx, appID, ownerID)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *AppService) GetApps(ctx context.Context, ownerID int) ([]models.App, error) {
	apps, err := a.appRepository.GetApps(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (a *AppService) DeleteApp(ctx context.Context, appID string, ownerID int) error {
	err := a.appRepository.DeleteApp(ctx, appID, ownerID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppService) UpdateApp(ctx context.Context, appID string, app DTO.UpdateApp, ownerID int) error {
	err := a.appRepository.UpdateApp(ctx, appID, app, ownerID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppService) GetAppStatus(ctx context.Context, appID string, ownerID int) (DTO.AppStatus, error) {
	return a.appStatusService.GetAppStatus(ctx, appID, ownerID)
}

func (a *AppService) CheckAppsStatus(ctx context.Context) ([]DTO.AppStatus, error) {
	return a.appStatusService.CheckAppsStatus(ctx)
}

func (a *AppService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	return a.appNotificationsService.SendNotifications(ctx, appsStatuses)
}

func (a *AppService) CheckRoutesStatus(ctx context.Context) error {
	return a.routeStatusService.CheckRoutesStatus(ctx)
}
