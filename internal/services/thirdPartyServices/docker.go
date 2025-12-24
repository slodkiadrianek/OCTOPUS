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

func (dc *DockerService) PauseContainer(ctx context.Context, appID string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	err = cli.ContainerPause(ctx, appID)

	return err
}

func (dc *DockerService) RestartContainer(ctx context.Context, appID string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, appID, containerTypes.StopOptions{})
	if err != nil {
		return err
	}
	err = cli.ContainerStart(ctx, appID, containerTypes.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (dc *DockerService) StartContainer(ctx context.Context, appID string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStart(ctx, appID, containerTypes.StartOptions{})

	return err
}

func (dc *DockerService) UnpauseContainer(ctx context.Context, appID string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerUnpause(ctx, appID)

	return err
}

func (dc *DockerService) StopContainer(ctx context.Context, appID string) error {
	cli, err := client.NewClientWithOpts(client.WithHost(dc.dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, appID, containerTypes.StopOptions{})

	return err
}

func (dc *DockerService) prepareContainersDataToInsert(containers []containerTypes.Summary, ownerID int, importedApps []models.App) []DTO.App {
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
					preparedIPAddress := strings.Split(splittedDockerHost, ":")[0]
					appsChan <- *DTO.NewApp(job.ID, preparedAppName, "", true, ownerID, preparedIPAddress, preparedPort)
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

func (dc *DockerService) ImportContainers(ctx context.Context, ownerID int) error {
	importedApps, err := dc.appRepository.GetApps(ctx, ownerID)
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
		dc.logger.Error("failed to list containers", err)
		return err
	}

	appsToInsert := dc.prepareContainersDataToInsert(containers, ownerID, importedApps)
	if len(appsToInsert) == 0 {
		return nil
	}
	err = dc.appRepository.InsertApp(ctx, appsToInsert)
	if err != nil {
		return err
	}

	return nil
}
