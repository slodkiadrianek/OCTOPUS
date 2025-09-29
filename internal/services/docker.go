package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type DockerService struct {
	DockerRepository *repository.DockerRepository
	AppRepository    *repository.AppRepository
	DockerHost       string
	Logger           *utils.Logger
}

func NewDockerService(dockerRepository *repository.DockerRepository, appRepository *repository.AppRepository, logger *utils.Logger, dockerHost string) *DockerService {
	return &DockerService{
		DockerRepository: dockerRepository,
		AppRepository:    appRepository,
		DockerHost:       dockerHost,
		Logger:           logger,
	}
}

func (dc *DockerService) PauseContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerPause(ctx, appId)
	return err
}

func (dc *DockerService) RestartContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerStop(ctx, appId, containertypes.StopOptions{})
	err = cli.ContainerStart(ctx, appId, containertypes.StartOptions{})
	return err
}

func (dc *DockerService) StartContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerStart(ctx, appId, containertypes.StartOptions{})
	return err
}

func (dc *DockerService) UnpauseContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerUnpause(ctx, appId)
	return err
}
func (dc *DockerService) StopContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerStop(ctx, appId, containertypes.StopOptions{})
	return err
}

func (dc *DockerService) ImportContainers(ctx context.Context, ownerId int) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		dc.Logger.Error("Failed to list containers", err)
		return err
	}
	//appsData := {}
	var appsData []DTO.App
	workerCount := runtime.NumCPU()
	jobs := make(chan containertypes.Summary, len(containers))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				preparedName := job.Names[0][1:]
				appsData = append(appsData, *DTO.NewApp(job.ID, preparedName, "", true, ownerId, "", ""))
			}
		}(i + 1)
	}
	for _, container := range containers {
		jobs <- container
	}
	close(jobs)
	wg.Wait()
	fmt.Println("Using", workerCount, "workers")
	err = dc.AppRepository.InsertApp(ctx, appsData)
	return err
}
