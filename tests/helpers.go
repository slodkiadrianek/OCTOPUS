package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	image2 "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

func CreateCacheService(loggerService *utils.Logger) interfaces.CacheService {

	cfg, err := config.SetConfig("../../../.env")
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

func CreateLogger() *utils.Logger {
	loggerService := utils.NewLogger("../../../logs", "2006-01-02 15:04:05")
	loggerService.InitializeLogger()
	return loggerService
}

func Ptr[T ~int | ~string](v T) *T {
	return &v
}

func CreateTestContainer(image string, cmd []string, loggerService *utils.Logger, dockerHost string) (string, error) {
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

	dec := json.NewDecoder(out)
	var status map[string]interface{}
	for dec.More() {
		if err := dec.Decode(&status); err == nil {
			if progress, ok := status["status"]; ok {
				fmt.Println(progress)
			}
		}
	}
	_, err = io.Copy(io.Discard, out)
	if err != nil {
		return "", err
	}

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

func KillAndRemoveContainer(ctx context.Context, containerID string, loggerService *utils.Logger,
	dockerHost string) error {
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
