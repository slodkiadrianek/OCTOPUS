package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

//func TestAppService_CheckAppsStatus(t *testing.T) {
//	env, err := config.SetConfig("../../.env")
//	if err != nil {
//		panic(err)
//	}
//	type args struct {
//		name          string
//		expectedError *string
//		dockerHost string
//		setupMock     func() (appRepository, CacheService)
//	}
//	tests := []args{
//		{
//			name:          "Proper data with data didn't save  in cache",
//			expectedError: nil,
//			setupMock: func() (appRepository, CacheService) {
//				mCache := new(mocks.MockCacheService)
//				mApp := new(mocks.MockAppRepository)
//				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(0), nil)
//				mApp.On("GetAppStatus", mock.Anything, mock.Anything, mock.Anything).Return(
//					DTO.AppStatus{AppId: "23r32"}, nil)
//				return mApp, mCache
//			},
//		},
//		{
//			name:          "Proper data with data saved in cache",
//			expectedError: nil,
//			setupMock: func() (appRepository, CacheService) {
//				mCache := new(mocks.MockCacheService)
//				mApp := new(mocks.MockAppRepository)
//				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
//				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`{
//				  "app_id": "com.example.myapp.prod",
//				  "status": "RUNNING",
//				  "changed_at": "2025-10-08T14:30:00Z",
//				  "duration": 7800000000000
//				}`, nil)
//				return mApp, mCache
//			},
//		},
//		{
//			name:          "Failed to get data from cache",
//			expectedError: ptr("Internal server error"),
//			setupMock: func() (appRepository, CacheService) {
//				mCache := new(mocks.MockCacheService)
//				mApp := new(mocks.MockAppRepository)
//				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
//				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`
//				`, errors.New("Failed to get data"))
//				return mApp, mCache
//			},
//		},
//		{
//			name:          "Wrong data format provided from cache",
//			expectedError: ptr("Internal server error"),
//			setupMock: func() (appRepository, CacheService) {
//				mCache := new(mocks.MockCacheService)
//				mApp := new(mocks.MockAppRepository)
//				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(1), nil)
//				mCache.On("GetData", mock.Anything, "status-123e23e23").Return(`
//				invalid-format`, nil)
//				return mApp, mCache
//			},
//		},
//		{
//			name:          "Failed to get data from database",
//			expectedError: ptr("Failed to get data from db"),
//			setupMock: func() (appRepository, CacheService) {
//				mCache := new(mocks.MockCacheService)
//				mApp := new(mocks.MockAppRepository)
//				mCache.On("ExistsData", mock.Anything, "status-123e23e23").Return(int64(0), nil)
//				mApp.On("GetAppStatus", mock.Anything, mock.Anything, mock.Anything).Return(
//					DTO.AppStatus{}, errors.New("Failed to get data from db"))
//				return mApp, mCache
//			},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			ctx := context.Background()
//			loggerService := createLogger()
//			appRepository, cacheService := test.setupMock()
//			appService := NewAppService(appRepository, loggerService, cacheService, env.DockerHost)
//			app, err := appService.GetAppStatus(ctx,
//				"123e23e23", 543)
//			fmt.Println(app, err)
//			if test.expectedError == nil {
//				assert.NoError(t, err)
//				assert.NotEmpty(t, app)
//			} else {
//				assert.Error(t, err)
//				assert.Empty(t, app)
//				assert.Contains(t, err.Error(), *test.expectedError)
//			}
//		})
//	}
//}
