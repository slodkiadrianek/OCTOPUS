package server

import (
	"context"
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/stretchr/testify/mock"

	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServerService_GetServerMetrics(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		expectedError *string
		setupMock     func() interfaces.CacheService
	}

	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() interfaces.CacheService {
				m := new(mocks.MockCacheService)
				m.On("GetData", mock.Anything, "server_metrics").
					Return(`[{"CPU": 50, "Memory": 1024, "Disk":30,"Date": "2023-05-12T10:30:00Z"}]`, nil)
				return m
			},
		},
		{
			name:          "Error in cache service",
			expectedError: tests.Ptr("failed to get data"),
			setupMock: func() interfaces.CacheService {
				m := new(mocks.MockCacheService)
				m.On("GetData", mock.Anything, "server_metrics").
					Return("", errors.New("failed to get data"))
				return m
			},
		},
		{
			name:          "Error in unmarshal",
			expectedError: tests.Ptr("invalid"), // or "invalid character"
			setupMock: func() interfaces.CacheService {
				m := new(mocks.MockCacheService)
				m.On("GetData", mock.Anything, "server_metrics").
					Return("not-a-valid-json", nil)
				return m
			},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			cache := testScenario.setupMock()
			serverService := NewServerService(loggerService, cache)
			ctx := context.Background()
			res, err := serverService.GetServerMetrics(ctx)
			if testScenario.expectedError == nil {
				assert.NotEmpty(t, res)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, res)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *testScenario.expectedError)
			}
		})
	}
}

func TestServerService_InsertServerMetrics(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		expectedError *string
		setupMock     func() interfaces.CacheService
	}
	testsScenarios := []args{
		{
			name:          "Proper data with data saved in cache",
			expectedError: nil,
			setupMock: func() interfaces.CacheService {
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
			setupMock: func() interfaces.CacheService {
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
			expectedError: tests.Ptr("issue with cache service"),
			setupMock: func() interfaces.CacheService {
				m := new(mocks.MockCacheService)
				m.On("ExistsData", mock.Anything, "server_metrics").
					Return(int64(0), errors.New("issue with cache service"))
				return m
			},
		},
		{
			name:          "Problem with GetData",
			expectedError: tests.Ptr("issue with cache service"),
			setupMock: func() interfaces.CacheService {
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
			expectedError: tests.Ptr("issue with cache service"),
			setupMock: func() interfaces.CacheService {
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
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			cache := testScenario.setupMock()
			serverService := NewServerService(loggerService, cache)
			ctx := context.Background()
			err := serverService.InsertServerMetrics(ctx)
			if testScenario.expectedError == nil {
				assert.Nil(t, err)
			}
			if err != nil {
				assert.Equal(t, *testScenario.expectedError, err.Error())
			}
		})
	}
}

func TestServerService_GetServerInfo(t *testing.T) {
	loggerService := tests.CreateLogger()
	cacheService := tests.CreateCacheService(loggerService)
	serverService := NewServerService(loggerService, cacheService)
	res, err := serverService.GetServerInfo()
	if err == nil {
		assert.NotEmpty(t, res)
	} else {
		assert.Empty(t, res)
	}
}
