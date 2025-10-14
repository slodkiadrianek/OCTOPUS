package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppService_CreateApp(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("InsertApp", mock.Anything, mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to insert an app",
			expectedError: ptr("failed to insert an app"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("InsertApp", mock.Anything, mock.Anything).Return(errors.New("failed to insert an app"))
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			app := DTO.CreateApp{
				Name:        "test",
				Description: "",
				IpAddress:   "192.168.2.22",
				Port:        "3020",
			}
			err := appService.CreateApp(ctx, app, 345)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}

func TestAppService_GetApp(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApp", mock.Anything, mock.Anything, mock.Anything).Return(&models.App{
					Id: "ewfw4f",
				}, nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to get an app",
			expectedError: ptr("failed to get an app"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApp", mock.Anything, mock.Anything, mock.Anything).Return(&models.App{}, errors.New("failed to get an app"))
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			app, err := appService.GetApp(ctx, "hf9hrepuihfefui", 32)
			if test.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, app)
			} else {
				assert.Error(t, err)
				assert.Empty(t, app)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}

func TestAppService_GetApps(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApps", mock.Anything, mock.Anything).Return([]models.App{{
					Id: "ewfw4f",
				}}, nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to get an apps",
			expectedError: ptr("failed to get an apps"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApps", mock.Anything, mock.Anything).Return([]models.App{}, errors.New("failed to get an apps"))
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			app, err := appService.GetApps(ctx,
				32)
			if test.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, app)
			} else {
				assert.Error(t, err)
				assert.Empty(t, app)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}
func TestAppService_GetAppStatus(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Proper data with data didn't save  in cache",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(0), nil)
				mApp.On("GetAppStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					DTO.AppStatus{AppId: "23r32"}, nil)
				return mApp, mCache
			},
		},
		{
			name:          "Proper data with data saved in cache",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`{
				  "app_id": "com.example.myapp.prod",
				  "status": "RUNNING",
				  "changed_at": "2025-10-08T14:30:00Z",
				  "duration": 7800000000000
				}`, nil)
				return mApp, mCache
			},
		},
		{
			name:          "Failed to get data from cache",
			expectedError: ptr("Internal server error"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`
				`, errors.New("Failed to get data"))
				return mApp, mCache
			},
		},
		{
			name:          "Wrong data format provided from cache",
			expectedError: ptr("Internal server error"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`
				invalid-format`, nil)
				return mApp, mCache
			},
		},
		{
			name:          "Failed to get data from database",
			expectedError: ptr("Failed to get data from db"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(0), nil)
				mApp.On("GetAppStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					DTO.AppStatus{}, errors.New("Failed to get data from db"))
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			app, err := appService.GetAppStatus(ctx,
				"123e23e23", 543)
			fmt.Println(app, err)
			if test.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, app)
			} else {
				assert.Error(t, err)
				assert.Empty(t, app)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}

func TestAppService_DeleteApp(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("DeleteApp", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to delete an app",
			expectedError: ptr("Failed to delete an app"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("DeleteApp", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Failed to delete an app"))
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			err := appService.DeleteApp(ctx,
				"delete", 21)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}
func TestAppService_UpdateApp(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("UpdateApp", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to delete an app",
			expectedError: ptr("Failed to update an app"),
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("UpdateApp", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Failed to update an app"))
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			app := DTO.UpdateApp{Name: "Test", Description: "test", Port: "3020", IpAddress: "192.168.20.10"}
			err := appService.UpdateApp(ctx,
				"delete", app, 21)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}

func TestAppService_CheckAppsStatus(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		dockerHost    string
		setupMock     func(appId string) (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "Failed to get app to check",
			expectedError: ptr("failed to get app to check"),
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{},
					errors.New("failed to get app to check"))
				return mApp, mCache
			},
		},
		{
			name:          "Wrong docker host provided",
			expectedError: ptr("unable to parse docker host"),
			dockerHost:    "192.168.0.100",
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{},
					nil)
				return mApp, mCache
			},
		},
		{
			name:          "Successfully to inspected the container",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						Id:       appId,
						IsDocker: true,
					},
				},
					nil)
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(nil)
				mCache.On("SetData", mock.Anything, "status-"+appId, mock.Anything,
					mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "Failed to inspect container",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						Id:       "r32r23r",
						IsDocker: true,
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "Check not a container app",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						Id:        "r32r23r",
						IsDocker:  false,
						IpAddress: "192.168.0.100",
						Port:      env.Port,
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(nil)
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(nil)

				return mApp, mCache
			},
		},
		{
			name:          "failed to save status in cache",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						Id:        "r32r23r",
						IsDocker:  false,
						IpAddress: "192.168.0.100",
						Port:      env.Port,
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(errors.New("Failed to save app status in cache"))
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(nil)

				return mApp, mCache
			},
		},
		{
			name:          "App stopped",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						Id:        "r32r23r",
						IsDocker:  false,
						IpAddress: "192.168.0.100",
						Port:      "9999",
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(nil)
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(nil)

				return mApp, mCache
			},
		},
		{
			name:          "Failed to insert appStatuses to db",
			expectedError: ptr("Failed to insert appStatuses to db"),
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						Id:        "r32r23r",
						IsDocker:  false,
						IpAddress: "192.168.0.100",
						Port:      env.Port,
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(nil)
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(errors.New("Failed to insert appStatuses to db"))

				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			containerId, _ := createTestContainer("alpine", []string{"sleep", "10"},
				loggerService,
				env.DockerHost)
			appId := containerId
			appRepository, cacheService := test.setupMock(appId)
			appService := NewAppService(appRepository, loggerService, cacheService, test.dockerHost)
			_, err := appService.CheckAppsStatus(ctx)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
			err = killAndRemoveContainer(ctx, appId, loggerService, env.DockerHost)
			if err != nil {
				panic(err)
			}

		})
	}
}
func TestAppService_SendNotifications(t *testing.T) {
	env, err := config.SetConfig("../../.env")
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError *string
		appsStatuses  []DTO.AppStatus
		setupMock     func() (appRepository, CacheService)
	}
	tests := []args{
		{
			name:          "No app statuses",
			expectedError: nil,
			appsStatuses:  []DTO.AppStatus{},
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				return mApp, mCache
			},
		},
		{
			name:          "Failed to get users to send notifications",
			expectedError: ptr("Failed to get users to send notifications"),
			appsStatuses:  []DTO.AppStatus{{AppId: "32"}},
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetUsersToSendNotifications", mock.Anything,
					mock.Anything).Return([]models.SendNotificationTo{}, errors.New("Failed to get users to send notifications"))
				return mApp, mCache
			},
		},
		{
			name:          "Proper data",
			expectedError: nil,
			appsStatuses:  []DTO.AppStatus{{AppId: "32", Status: "running"}},
			setupMock: func() (appRepository, CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetUsersToSendNotifications", mock.Anything,
					mock.Anything).Return([]models.SendNotificationTo{
					{
						Status:               "running",
						SlackNotifications:   true,
						DiscordNotifications: true,
						DiscordWebhook:       "",
						SlackWebhook:         "https://hooks.slack.com/services/T026/B09AY4T/zunu2tPqHARDJ",
					},
				}, nil)
				return mApp, mCache
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			appRepository, cacheService := test.setupMock()
			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
			err := appService.SendNotifications(ctx,
				test.appsStatuses)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}
