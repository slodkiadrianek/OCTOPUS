package servicesApp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppService_GetAppStatus(t *testing.T) {
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
			name:          "Proper data with data didn't save  in cache",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(0), nil)
				mApp.On("GetAppStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					DTO.AppStatus{AppID: "23r32"}, nil)
				return mApp, mCache
			},
		},
		{
			name:          "Proper data with data saved in cache",
			expectedError: nil,
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
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
			name:          "failed to get data from cache",
			expectedError: errors.New("Internal server error"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`
				`, errors.New("failed to get data"))
				return mApp, mCache
			},
		},
		{
			name:          "Wrong data format provided from cache",
			expectedError: errors.New("Internal server error"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`
				invalid-format`, nil)
				return mApp, mCache
			},
		},
		{
			name:          "failed to get data from database",
			expectedError: errors.New("failed to get data from db"),
			setupMock: func() (interfaces.AppRepository, interfaces.CacheService) {
				mCache := new(mocks.MockCacheService)
				mApp := new(mocks.MockAppRepository)
				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(0), nil)
				mApp.On("GetAppStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					DTO.AppStatus{}, errors.New("failed to get data from db"))
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
			app, err := appService.GetAppStatus(ctx,
				"123e23e23", 543)
			fmt.Println(app, err)
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
