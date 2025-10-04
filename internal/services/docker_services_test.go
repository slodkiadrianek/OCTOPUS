package services

import (
	"context"
	"errors"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDockerService_ImportContainers(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		dockerHost    string
		expectedError *string
		setupMock     func() appRepository
	}
	tests := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			dockerHost:    "tcp://100.100.188.29:2375",
			setupMock: func() appRepository {
				m := new(mocks.MockAppRepository)
				m.On("InsertApp", mock.Anything, mock.Anything).Return(nil)
				return m
			},
		},
		{
			name:          "Invalid docker host",
			expectedError: ptr("no such host"),
			dockerHost:    "tcp://100.100.188.329:2375",
			setupMock: func() appRepository {
				m := new(mocks.MockAppRepository)
				m.On("InsertApp", mock.Anything, mock.Anything).Return(nil)
				return m
			},
		},
		{
			name:          "Insert to db error",
			expectedError: ptr("failed to insert data to db"),
			dockerHost:    "tcp://100.100.188.29:2375",
			setupMock: func() appRepository {
				m := new(mocks.MockAppRepository)
				m.On("InsertApp", mock.Anything, mock.Anything).Return(errors.New("failed to insert data to db"))
				return m
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appRepositoryMock := test.setupMock()
			dockerService := NewDockerService(appRepositoryMock, loggerService, test.dockerHost)
			ctx := context.Background()
			err := dockerService.ImportContainers(ctx, 1)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
		})
	}
}
