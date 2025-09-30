package services

import (
	"context"
	"runtime"
	"strings"
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

func (dc *DockerService) ImportContainers(ctx context.Context, ownerId int) error {
	cli, err := client.NewClientWithOpts(
		client.WithHost(dc.DockerHost),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		dc.Logger.Error("Failed to list containers", err)
		return err
	}

	workerCount := runtime.NumCPU()
	jobs := make(chan containertypes.Summary, len(containers))
	results := make(chan DTO.App, len(containers))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				if len(job.Names) == 0 {
					continue
				}
				preparedName := strings.TrimPrefix(job.Names[0], "/")
				app := DTO.NewApp(job.ID, preparedName, "", true, ownerId, "", "")
				results <- *app
			}
		}(i + 1)
	}

	for _, container := range containers {
		jobs <- container
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var appsData []DTO.App
	for app := range results {
		appsData = append(appsData, app)
	}

	dc.Logger.Info("Collected apps", "count", len(appsData), "workers", workerCount)

	if err := dc.AppRepository.InsertApp(ctx, appsData); err != nil {
		dc.Logger.Error("Failed to insert apps", err)
		return err
	}

	return nil
}
