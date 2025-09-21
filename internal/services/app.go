package services

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type AppService struct {
	AppRepository *repository.AppRepository
	Logger        *logger.Logger
	CacheService  *config.CacheService
	DockerHost    string
}

func NewAppService(appRepository *repository.AppRepository, logger *logger.Logger, cacheService *config.CacheService, dockerHost string) *AppService {
	return &AppService{
		AppRepository: appRepository,
		Logger:        logger,
		CacheService:  cacheService,
		DockerHost:    dockerHost,
	}
}

func (a *AppService) CreateApp(ctx context.Context, app schema.CreateApp, ownerId int) error {
	id, err := utils.GenerateID()
	if err != nil {
		return err
	}
	appDto := DTO.NewApp(id, app.Name, app.Description, false, ownerId, "", "")
	err = a.AppRepository.InsertApp(ctx, []DTO.App{*appDto})
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetApp(ctx context.Context, id int) (*models.App, error) {
	app, err := a.AppRepository.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) UpdateApp(ctx context.Context, id int, app schema.UpdateApp) error {
	// appDto := DTO.NewUpdateApp(id, app.Name, app.Description, app.DbLink, app.ApiLink, app.DiscordWebhook, app.SlackWebhook)
	return nil
}

func (a *AppService) DeleteApp(ctx context.Context, id int) error {
	err := a.AppRepository.DeleteApp(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) CheckAppsStatus(ctx context.Context) error {
	apps, err := a.AppRepository.GetAppsToCheck(ctx)
	if err != nil {
		return err
	}
	workerCount := runtime.NumCPU()
	cli, err := client.NewClientWithOpts(client.WithHost(a.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	var appsStatuses []DTO.AppStatus
	defer cli.Close()
	jobs := make(chan *models.AppToCheck, len(apps))
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for job := range jobs {
				if job.IsDocker {
					container, err := cli.ContainerInspect(ctx, job.Id)
					if err != nil {
						a.Logger.Error("Failed to check status inside of a container", err)
						continue
					}
					status := container.State.Status
					startedTime, err := time.Parse(time.RFC3339, container.State.StartedAt)
					if err != nil {
						a.Logger.Error("Failed to parse time", err)
						continue
					}
					duration := time.Since(startedTime)
					changedAtTime, err := time.Parse(time.RFC3339, container.State.StartedAt)
					if err != nil {
						a.Logger.Error("Failed to parse time", err)
						continue
					}
					appStatus := *DTO.NewAppStatus(job.Id, status, changedAtTime, duration)
					appsStatuses = append(appsStatuses, appStatus)
					bodyBytes, err := utils.MarshalData(appStatus)
					if err != nil {
						a.Logger.Error("Failed to convert data to json", map[string]any{
							"data":  appStatus,
							"error": err.Error(),
						})
						continue
					}
					err = a.CacheService.SetData(ctx, "status-"+job.Id, string(bodyBytes), time.Minute*2)
					if err != nil {
						a.Logger.Error("Failed to setData in cache", map[string]any{
							"data":  appStatus,
							"error": err.Error(),
						})
						continue
					}
				} else {
					address := fmt.Sprintf("%s:%d", job.IpAddress, job.Port)
					conn, err := net.DialTimeout("tcp", address, time.Second*3)
					if err != nil {
						a.Logger.Error("Failed to check status inside of a container", err)
						continue
					}
					defer conn.Close()
					status := "running"
					duration := time.Since(time.Now())
					changedAt := time.Now()
					appsStatuses = append(appsStatuses, *DTO.NewAppStatus(job.Id, status, changedAt, duration))
				}
			}
		}(i + 1)
	}
	for _, app := range apps {
		jobs <- app
	}
	close(jobs)
	wg.Wait()
	err = a.AppRepository.InsertAppStatuses(ctx, appsStatuses)
	return err
}
