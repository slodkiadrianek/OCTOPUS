package services

import (
	"context"
	"testing"

	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/mock"
)

func TestAppService_CreateApp(t *testing.T) {
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := createLogger()
			cacheService := ca
		})
	}
}
