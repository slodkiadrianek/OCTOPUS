package thirdPartyServices

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"

	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
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

func (dc *DockerService) prepareContainersDataToInsert(containers []containerTypes.Summary, ownerId int, importedApps []models.App) []DTO.App {
	workerCount := runtime.NumCPU()
	jobs := make(chan containerTypes.Summary, len(containers))
	appsChan := make(chan DTO.App, len(containers))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		outer:
			for job := range jobs {
				for _, app := range importedApps {
					if app.ID == job.ID {
						break outer
					}
				}
				preparedAppName := job.Names[0][1:]
				if len(job.Ports) > 0 {
					preparedPort := fmt.Sprintf("%d", job.Ports[0].PrivatePort)
					splittedDockerHost := strings.Split(dc.dockerHost, "//")[1]
					preparedIpAddress := strings.Split(splittedDockerHost, ":")[0]
					appsChan <- *DTO.NewApp(job.ID, preparedAppName, "", true, ownerId, preparedIpAddress, preparedPort)
				}
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
	importedApps, err := dc.appRepository.GetApps(ctx, ownerId)
	if err != nil {
		return err
	}

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

	appsToInsert := dc.prepareContainersDataToInsert(containers, ownerId, importedApps)
	if len(appsToInsert) == 0 {
		return nil
	}
	err = dc.appRepository.InsertApp(ctx, appsToInsert)
	if err != nil {
		return err
	}

	return nil
}
