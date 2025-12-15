package thirdPartyServices

import (
	"context"
	"runtime"
	"sync"

	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type DockerService struct {
	appRepository interfaces.AppRepository
	dockerHost    string
	logger        utils.LoggerService
}

func NewDockerService(appRepository interfaces.AppRepository, logger utils.LoggerService,
	dockerHost string,
) *DockerService {
	return &DockerService{
		appRepository: appRepository,
		dockerHost:    dockerHost,
		logger:        logger,
	}
}

func (dc *DockerService) PauseContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	err = cli.ContainerPause(ctx, appId)

	return err
}

func (dc *DockerService) RestartContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, appId, containerTypes.StopOptions{})

	err = cli.ContainerStart(ctx, appId, containerTypes.StartOptions{})

	return err
}

func (dc *DockerService) StartContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStart(ctx, appId, containerTypes.StartOptions{})

	return err
}

func (dc *DockerService) UnpauseContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerUnpause(ctx, appId)

	return err
}

func (dc *DockerService) StopContainer(ctx context.Context, appId string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, appId, containerTypes.StopOptions{})

	return err
}
func (dc *DockerService) prepareContainersDataToInert(containers []containerTypes.Summary, ownerId int) []DTO.App {
	workerCount := runtime.NumCPU()
	jobs := make(chan containerTypes.Summary, len(containers))
	appsChan := make(chan DTO.App, len(containers))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				preparedAppName := job.Names[0][1:]
				appsChan <- *DTO.NewApp(job.ID, preparedAppName, "", true, ownerId, "", "")
			}
		}()
	}
	for _, container := range containers {
		jobs <- container
	}
	close(jobs)
	wg.Wait()
	close(appsChan)
	appsToInsert := make([]DTO.App, 0, len(containers))
	for app := range appsChan {
		appsToInsert = append(appsToInsert, app)
	}

	return appsToInsert
}

func (dc *DockerService) ImportContainers(ctx context.Context, ownerId int) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, containerTypes.ListOptions{})
	if err != nil {
		dc.logger.Error("Failed to list containers", err)
		return err
	}

	appsToInsert := dc.prepareContainersDataToInert(containers, ownerId)
	err = dc.appRepository.InsertApp(ctx, appsToInsert)
	if err != nil {
		return err
	}

	return nil
}
