package thirdPartyServices

import (
	"context"
	"errors"
	"time"

	"testing"

	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDockerService_ImportContainers(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		dockerHost    string
		expectedError error
		setupMock     func() interfaces.AppRepository
	}
	testsScenarios := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			dockerHost:    "tcp://100.100.188.29:2375",
			setupMock: func() interfaces.AppRepository {
				m := new(mocks.MockAppRepository)
				m.On("InsertApp", mock.Anything, mock.Anything).Return(nil)
				return m
			},
		},
		{
			name:          "Invalid docker host",
			expectedError: errors.New("unable to parse docker host"),
			dockerHost:    "",
			setupMock: func() interfaces.AppRepository {
				m := new(mocks.MockAppRepository)
				m.On("InsertApp", mock.Anything, mock.Anything).Return(nil)
				return m
			},
		},
		{
			name:          "Insert to db error",
			expectedError: errors.New("failed to insert data to db"),
			dockerHost:    "tcp://100.100.188.29:2375",
			setupMock: func() interfaces.AppRepository {
				m := new(mocks.MockAppRepository)
				m.On("InsertApp", mock.Anything, mock.Anything).Return(errors.New("failed to insert data to db"))
				return m
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			appRepositoryMock := testScenario.setupMock()
			dockerService := NewDockerService(appRepositoryMock, loggerService, testScenario.dockerHost)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			err := dockerService.ImportContainers(ctx, 1)
			defer cancel()
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestDockerService_PauseContainer(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: errors.New("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app ID",
			expectedError: errors.New("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, testScenario.dockerHost)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel()
			var appId string
			if testScenario.appId == "" {
				containerId, _ := tests.CreateTestContainer("alpine", []string{"sleep", "20"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = testScenario.appId
			}
			err := dockerService.PauseContainer(ctx, appId)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			if testScenario.appId == "" {
				err := tests.KillAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestDockerService_RestartContainer(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: errors.New("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app ID",
			expectedError: errors.New("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, testScenario.dockerHost)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel()
			var appId string
			if testScenario.appId == "" {
				containerId, _ := tests.CreateTestContainer("alpine", []string{"sleep", "20"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = testScenario.appId
			}
			err := dockerService.RestartContainer(ctx, appId)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			if testScenario.appId == "" {
				err := tests.KillAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestDockerService_UnpauseContainer(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: errors.New("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app ID",
			expectedError: errors.New("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, testScenario.dockerHost)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel()
			var appId string
			if testScenario.appId == "" {
				containerId, _ := tests.CreateTestContainer("alpine", []string{"sleep", "20"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = testScenario.appId
			}
			err := dockerService.PauseContainer(ctx, appId)
			err = dockerService.UnpauseContainer(ctx, appId)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			if testScenario.appId == "" {
				err := tests.KillAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestDockerService_StartContainer(t *testing.T) {
	loggerService := tests.CreateLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: errors.New("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app ID",
			expectedError: errors.New("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, testScenario.dockerHost)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel()
			var appId string
			if testScenario.appId == "" {
				containerId, _ := tests.CreateTestContainer("alpine", []string{"sleep", "20"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = testScenario.appId
			}
			_ = dockerService.StopContainer(ctx, appId)
			err := dockerService.StartContainer(ctx, appId)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			if testScenario.appId == "" {
				err := tests.KillAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}
