package servicesApp

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppService_CreateApp(t *testing.T) {
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError error
		setupMock     func() (interfaces.AppRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("InsertApp", mock.Anything, mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to insert an app",
			expectedError: errors.New("failed to insert an app"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("InsertApp", mock.Anything, mock.Anything).Return(errors.New("failed to insert an app"))
				return mApp, mCache
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			appRepository, cacheService := testScenario.setupMock()
			routeRepository := repository.NewRouteRepository(&sql.DB{}, loggerService)
			appStatusService := NewAppStatusService(appRepository, cacheService, loggerService, env.DockerHost)
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			appService := NewAppService(appRepository, loggerService, appStatusService, appNotificationsService, routeStatusService)
			app := DTO.CreateApp{
				Name:        "test",
				Description: "",
				IPAddress:   "192.168.2.22",
				Port:        "3020",
			}
			err := appService.CreateApp(ctx, app, 345)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestAppService_GetApp(t *testing.T) {
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError error
		setupMock     func() (interfaces.AppRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApp", mock.Anything, mock.Anything, mock.Anything).Return(&models.App{
					ID: "ewfw4f",
				}, nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to get an app",
			expectedError: errors.New("failed to get an app"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApp", mock.Anything, mock.Anything, mock.Anything).Return(&models.App{}, errors.New("failed to get an app"))
				return mApp, mCache
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			appRepository, cacheService := testScenario.setupMock()
			routeRepository := repository.NewRouteRepository(&sql.DB{}, loggerService)
			appStatusService := NewAppStatusService(appRepository, cacheService, loggerService, env.DockerHost)
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			appService := NewAppService(appRepository, loggerService, appStatusService, appNotificationsService, routeStatusService)
			app, err := appService.GetApp(ctx, "hf9hrepuihfefui", 32)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, app)
			} else {
				assert.Error(t, err)
				assert.Empty(t, app)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestAppService_GetApps(t *testing.T) {
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError error
		setupMock     func() (interfaces.AppRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApps", mock.Anything, mock.Anything).Return([]models.App{{
					ID: "ewfw4f",
				}}, nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to get an apps",
			expectedError: errors.New("failed to get an apps"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetApps", mock.Anything, mock.Anything).Return([]models.App{}, errors.New("failed to get an apps"))
				return mApp, mCache
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			appRepository, cacheService := testScenario.setupMock()
			routeRepository := repository.NewRouteRepository(&sql.DB{}, loggerService)
			appStatusService := NewAppStatusService(appRepository, cacheService, loggerService, env.DockerHost)
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			appService := NewAppService(appRepository, loggerService, appStatusService, appNotificationsService, routeStatusService)
			app, err := appService.GetApps(ctx,
				32)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, app)
			} else {
				assert.Error(t, err)
				assert.Empty(t, app)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestAppService_DeleteApp(t *testing.T) {
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError error
		setupMock     func() (interfaces.AppRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("DeleteApp", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to delete an app",
			expectedError: errors.New("failed to delete an app"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("DeleteApp", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failed to delete an app"))
				return mApp, mCache
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			appRepository, cacheService := testScenario.setupMock()
			routeRepository := repository.NewRouteRepository(&sql.DB{}, loggerService)
			appStatusService := NewAppStatusService(appRepository, cacheService, loggerService, env.DockerHost)
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			appService := NewAppService(appRepository, loggerService, appStatusService, appNotificationsService, routeStatusService)
			err := appService.DeleteApp(ctx,
				"delete", 21)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestAppService_UpdateApp(t *testing.T) {
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError error
		setupMock     func() (interfaces.AppRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("UpdateApp", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to delete an app",
			expectedError: errors.New("failed to update an app"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("UpdateApp", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failed to update an app"))
				return mApp, mCache
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			appRepository, cacheService := testScenario.setupMock()
			routeRepository := repository.NewRouteRepository(&sql.DB{}, loggerService)
			appStatusService := NewAppStatusService(appRepository, cacheService, loggerService, env.DockerHost)
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			appService := NewAppService(appRepository, loggerService, appStatusService, appNotificationsService, routeStatusService)
			app := DTO.UpdateApp{Name: "Test", Description: "test", Port: "3020", IPAddress: "192.168.20.10"}
			err := appService.UpdateApp(ctx,
				"delete", app, 21)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestAppService_CheckAppsStatus(t *testing.T) {
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	type args struct {
		name          string
		expectedError error
		dockerHost    string
		setupMock     func(appId string) (interfaces.AppRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Failed to get app to check",
			expectedError: errors.New("failed to get app to check"),
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{},
					errors.New("failed to get app to check"))
				return mApp, mCache
			},
		},
		{
			name:          "Wrong docker host provided",
			expectedError: errors.New("unable to parse docker host"),
			dockerHost:    "192.168.0.100",
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
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
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						ID:       appId,
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
			name:          "failed to inspect container",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						ID:       "r32r23r",
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
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						ID:        "r32r23r",
						IsDocker:  false,
						IPAddress: "192.168.0.100",
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
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						ID:        "r32r23r",
						IsDocker:  false,
						IPAddress: "192.168.0.100",
						Port:      env.Port,
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(errors.New("failed to save app status in cache"))
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(nil)

				return mApp, mCache
			},
		},
		{
			name:          "App stopped",
			expectedError: nil,
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						ID:        "r32r23r",
						IsDocker:  false,
						IPAddress: "192.168.0.100",
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
			name:          "failed to insert appStatuses to db",
			expectedError: errors.New("failed to insert appStatuses to db"),
			dockerHost:    env.DockerHost,
			setupMock: func(appId string) (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetAppsToCheck", mock.Anything).Return([]*models.AppToCheck{
					{
						ID:        "r32r23r",
						IsDocker:  false,
						IPAddress: "192.168.0.100",
						Port:      env.Port,
					},
				},
					nil)
				mCache.On("SetData", mock.Anything, "status-r32r23r", mock.Anything,
					mock.Anything).Return(nil)
				mApp.On("InsertAppStatuses", mock.Anything, mock.Anything).Return(errors.New("failed to insert appStatuses to db"))

				return mApp, mCache
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			containerId, _ := tests.CreateTestContainer("alpine", []string{"sleep", "10"},
				loggerService,
				env.DockerHost)
			appId := containerId
			appRepository, cacheService := testScenario.setupMock(appId)
			routeRepository := repository.NewRouteRepository(&sql.DB{}, loggerService)
			appStatusService := NewAppStatusService(appRepository, cacheService, loggerService, testScenario.dockerHost)
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			appService := NewAppService(appRepository, loggerService, appStatusService, appNotificationsService, routeStatusService)
			_, err := appService.CheckAppsStatus(ctx)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			err = tests.KillAndRemoveContainer(ctx, appId, loggerService, env.DockerHost)
			if err != nil {
				panic(err)
			}
		})
	}
}
