package interfaces

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
)

type AppRepository interface {
	InsertApp(ctx context.Context, app []DTO.App) error
	GetApp(ctx context.Context, id string, ownerID int) (*models.App, error)
	GetApps(ctx context.Context, ownerID int) ([]models.App, error)
	DeleteApp(ctx context.Context, id string, ownerID int) error
	GetAppStatus(ctx context.Context, id string, ownerID int) (DTO.AppStatus, error)
	GetAppsToCheck(ctx context.Context) ([]*models.AppToCheck, error)
	UpdateApp(ctx context.Context, appID string, app DTO.UpdateApp, ownerID int) error
	InsertAppStatuses(ctx context.Context, appsStatuses []DTO.AppStatus) error
	GetUsersToSendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) ([]models.NotificationInfo, error)
}

type AppNotificationsService interface {
	SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error
}
type AppStatusService interface {
	GetAppStatus(ctx context.Context, appID string, ownerID int) (DTO.AppStatus, error)
	CheckAppsStatus(ctx context.Context) ([]DTO.AppStatus, error)
}
