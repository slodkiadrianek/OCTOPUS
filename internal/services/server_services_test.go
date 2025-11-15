package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func createCacheService(loggerService *utils.Logger) CacheService {

	cfg, err := config.SetConfig("../../.env")
	if err != nil {
		loggerService.Error("Failed to load config", err)
		return nil
	}
	err = cfg.Validate()
	if err != nil {
		loggerService.Error("Configuration validation failed", err)
		return nil
	}
	cacheService, err := config.NewCacheService(cfg.CacheLink)
	if err != nil {
		loggerService.Error("Failed to connect to cache", err)
		return nil
	}
	return cacheService
}

func createLogger() *utils.Logger {
	loggerService := utils.NewLogger("../../logs", "2006-01-02 15:04:05")
	loggerService.InitializeLogger()
	return loggerService
}
func ptr(s string) *string {
	return &s
}

func TestServerService_GetServerMetrics(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		expectedError *string
		setupMock     func() CacheService
	}

	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("GetData", mock.Anything, "server_metrics").
					Return(`[{"CPU": 50, "Memory": 1024, "Disk":30,"Date": "2023-05-12T10:30:00Z"}]`, nil)
				return m
			},
		},
		{
			name:          "Error in cache service",
			expectedError: ptr("failed to get data"),
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("GetData", mock.Anything, "server_metrics").
					Return("", errors.New("failed to get data"))
				return m
			},
		},
		{
			name:          "Error in unmarshal",
			expectedError: ptr("invalid"), // or "invalid character"
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("GetData", mock.Anything, "server_metrics").
					Return("not-a-valid-json", nil)
				return m
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := test.setupMock()
			serverService := NewServerService(loggerService, cache)
			ctx := context.Background()
			res, err := serverService.GetServerMetrics(ctx)
			if test.expectedError == nil {
				assert.NotEmpty(t, res)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, res)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}

func TestServerService_InsertServerMetrics(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		expectedError *string
		setupMock     func() CacheService
	}
	tests := []args{
		{
			name:          "Proper data with data saved in cache",
			expectedError: nil,
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("ExistsData", mock.Anything, "server_metrics").
					Return(int64(1), nil)
				m.On("GetData", mock.Anything, "server_metrics").
					Return(`[{"CPU": 50, "Memory": 1024, "Disk":30,"Date": "2023-05-12T10:30:00Z"}]`, nil)
				m.On("SetData", mock.Anything, "server_metrics", mock.Anything, mock.Anything).
					Return(nil)
				return m
			},
		},
		{
			name:          "Proper data with data saved in cache",
			expectedError: nil,
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("ExistsData", mock.Anything, "server_metrics").
					Return(int64(0), nil)
				m.On("SetData", mock.Anything, "server_metrics", mock.Anything, mock.Anything).
					Return(nil)
				return m
			},
		},
		{
			name:          "Problem with ExistsData",
			expectedError: ptr("issue with cache service"),
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("ExistsData", mock.Anything, "server_metrics").
					Return(int64(0), errors.New("issue with cache service"))
				return m
			},
		},
		{
			name:          "Problem with GetData",
			expectedError: ptr("issue with cache service"),
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("ExistsData", mock.Anything, "server_metrics").
					Return(int64(1), nil)
				m.On("GetData", mock.Anything, "server_metrics").
					Return("", errors.New("issue with cache service"))
				return m
			},
		},
		{
			name:          "Problem with GetData",
			expectedError: ptr("issue with cache service"),
			setupMock: func() CacheService {
				m := new(mocks.MockCacheService)
				m.On("ExistsData", mock.Anything, "server_metrics").
					Return(int64(1), nil)
				m.On("GetData", mock.Anything, "server_metrics").
					Return(`[{"CPU": 50, "Memory": 1024, "Disk":30,"Date": "2023-05-12T10:30:00Z"}]`, nil)
				m.On("SetData", mock.Anything, "server_metrics", mock.Anything, mock.Anything).
					Return(errors.New("issue with cache service"))
				return m
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := test.setupMock()
			serverService := NewServerService(loggerService, cache)
			ctx := context.Background()
			err := serverService.InsertServerMetrics(ctx)
			if test.expectedError == nil {
				assert.Nil(t, err)
			}
			if err != nil {
				assert.Equal(t, *test.expectedError, err.Error())
			}
		})
	}
}

func TestServerService_GetServerInfo(t *testing.T) {
	loggerService := createLogger()
	cacheService := createCacheService(loggerService)
	serverService := NewServerService(loggerService, cacheService)
	res, err := serverService.GetServerInfo()
	if err == nil {
		assert.NotEmpty(t, res)
	} else {
		assert.Empty(t, res)
	}
}
