package services

import (
	"context"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func CreateCacheService() CacheService {
	loggerService := utils.NewLogger("./logs", "2006-01-02 15:04:05")
	loggerService.CreateLogger()
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

func TestGetServerMetrics(t *testing.T) {
	MockedCacheService := mocks.NewMockCacheService()
	cacheService := CreateCacheService()
	//serverService :=
	type args struct {
		name          string
		expectedError *string
		cacheService  CacheService
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			cacheService:  cacheService,
		},
		{
			name:          "Error in cache service",
			expectedError: nil,
			cacheService:  MockedCacheService,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			res, err := serverService.GetServerMetrics(ctx)
			if test.expectedError == nil {
				assert.NotEmpty(t, res)
			}
			if err != nil {
				assert.Empty(t, res)
			}
		})
	}
}

func TestInsertServerMetrics(t *testing.T) {
	serverService := CreateServerService()
	type args struct {
		name          string
		expectedError *string
	}
	tests := []args{
		{
			name:          "Proper data",
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			err := serverService.InsertServerMetrics(ctx)
			if test.expectedError == nil {
				assert.Nil(t, err)
			}
			if err != nil {
				assert.Equal(t, *test.expectedError, err)
			}
		})
	}
}
