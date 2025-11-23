package mocks

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAppRepository struct {
	mock.Mock
}

func (m *MockAppRepository) InsertApp(ctx context.Context, app []DTO.App) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockAppRepository) GetApp(ctx context.Context, id string, ownerId int) (*models.App, error) {
	args := m.Called(ctx, id, ownerId)
	return args.Get(0).(*models.App), args.Error(1)
}

func (m *MockAppRepository) GetApps(ctx context.Context, ownerId int) ([]models.App, error) {
	args := m.Called(ctx, ownerId)
	return args.Get(0).([]models.App), args.Error(1)
}

func (m *MockAppRepository) DeleteApp(ctx context.Context, id string, ownerId int) error {
	args := m.Called(ctx, id, ownerId)
	return args.Error(0)
}

func (m *MockAppRepository) GetAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error) {
	args := m.Called(ctx, id, ownerId)
	return args.Get(0).(DTO.AppStatus), args.Error(1)
}

func (m *MockAppRepository) GetAppsToCheck(ctx context.Context) ([]*models.AppToCheck, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.AppToCheck), args.Error(1)
}

func (m *MockAppRepository) UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error {
	args := m.Called(ctx, appId, app, ownerId)
	return args.Error(0)
}

func (m *MockAppRepository) InsertAppStatuses(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	args := m.Called(ctx, appsStatuses)
	return args.Error(0)
}

func (m *MockAppRepository) GetUsersToSendNotifications(ctx context.Context,
	appsStatuses []DTO.AppStatus,
) ([]models.NotificationInfo, error) {
	args := m.Called(ctx, appsStatuses)
	return args.Get(0).([]models.NotificationInfo), args.Error(1)
}
