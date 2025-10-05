package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	image2 "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestContainer(image string, cmd []string, loggerService *utils.Logger, dockerHost string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithHost(dockerHost))
	if err != nil {
		return "", err
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, image, image2.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer out.Close()

	// Drain and display pull progress to make sure it's done
	dec := json.NewDecoder(out)
	var status map[string]interface{}
	for dec.More() {
		if err := dec.Decode(&status); err == nil {
			if progress, ok := status["status"]; ok {
				fmt.Println(progress)
			}
		}
	}
	io.Copy(io.Discard, out) // ensure it's fully read

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		loggerService.Error("failed to create a new container", err)
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		loggerService.Error("failed to create a start container", err)
		return "", err
	}

	return resp.ID, nil
}

func killAndRemoveContainer(ctx context.Context, containerID string, loggerService *utils.Logger, dockerHost string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dockerHost))
	if err != nil {
		return err
	}
	defer cli.Close()

	timeout := 5 * time.Second
	timeAsInt := int(timeout)
	loggerService.Info("Stopping container", containerID)
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeAsInt}); err != nil {
		loggerService.Info("ContainerStop failed, trying ContainerKill...", err)
		_ = cli.ContainerKill(ctx, containerID, "SIGKILL")
	}
	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		loggerService.Error("failed to remove container", err)
		return err
	}

	return nil
}
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
			expectedError: ptr("unable to parse docker host"),
			dockerHost:    "",
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

func TestDockerService_PauseContainer(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError *string
	}
	tests := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: ptr("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app Id",
			expectedError: ptr("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, test.dockerHost)
			ctx := context.Background()
			var appId string
			if test.appId == "" {
				containerId, _ := createTestContainer("alpine", []string{"sleep", "60"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = test.appId
			}
			err := dockerService.PauseContainer(ctx, appId)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
			if test.appId == "" {
				err := killAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestDockerService_RestartContainer(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError *string
	}
	tests := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: ptr("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app Id",
			expectedError: ptr("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, test.dockerHost)
			ctx := context.Background()
			var appId string
			if test.appId == "" {
				containerId, _ := createTestContainer("alpine", []string{"sleep", "60"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = test.appId
			}
			err := dockerService.RestartContainer(ctx, appId)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
			if test.appId == "" {
				err := killAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestDockerService_UnpauseContainer(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError *string
	}
	tests := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: ptr("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app Id",
			expectedError: ptr("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, test.dockerHost)
			ctx := context.Background()
			var appId string
			if test.appId == "" {
				containerId, _ := createTestContainer("alpine", []string{"sleep", "60"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = test.appId
			}
			err := dockerService.PauseContainer(ctx, appId)
			err = dockerService.UnpauseContainer(ctx, appId)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
			if test.appId == "" {
				err := killAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestDockerService_StartContainer(t *testing.T) {
	loggerService := createLogger()
	type args struct {
		name          string
		dockerHost    string
		appId         string
		expectedError *string
	}
	tests := []args{
		{
			name:          "Test with proper data",
			expectedError: nil,
			appId:         "",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
		{
			name:          "Invalid docker host",
			expectedError: ptr("unable to parse docker host"),
			appId:         "",
			dockerHost:    "",
		},
		{
			name:          "Invalid app Id",
			expectedError: ptr("No such container"),
			appId:         "e9530eae6aa752adf79b79a2d9c1398fe59eee4a3d786734d9e2076e62415772",
			dockerHost:    "tcp://100.100.188.29:2375",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appRepositoryMock := new(mocks.MockAppRepository)
			dockerService := NewDockerService(appRepositoryMock, loggerService, test.dockerHost)
			ctx := context.Background()
			var appId string
			if test.appId == "" {
				containerId, _ := createTestContainer("alpine", []string{"sleep", "60"},
					loggerService,
					"tcp://100.100.188.29:2375")
				appId = containerId
			} else {
				appId = test.appId
			}
			_ = dockerService.StopContainer(ctx, appId)
			err := dockerService.StartContainer(ctx, appId)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *test.expectedError)
			}
			if test.appId == "" {
				err := killAndRemoveContainer(ctx, appId, loggerService, "tcp://100.100.188.29:2375")
				if err != nil {
					panic(err)
				}
			}
		})
	}
}
