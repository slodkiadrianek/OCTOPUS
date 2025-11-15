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
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AppStatusService struct {
	AppRepository appRepository
	CacheService  CacheService
	LoggerService utils.Logger
	DockerHost    string
}

func NewAppStatusService(appRepository appRepository, cacheService CacheService, loggerService utils.Logger, dockerHost string) *AppStatusService {
	return &AppStatusService{
		AppRepository: appRepository,
		CacheService:  cacheService,
		LoggerService: loggerService,
		DockerHost:    dockerHost,
	}
}

func (as *AppStatusService) readAppStatusFromCache(ctx context.Context, cacheKey string) (DTO.AppStatus, error) {
	appStatusAsJson, err := as.CacheService.GetData(ctx, cacheKey)
	if err != nil {
		as.LoggerService.Warn("Failed to get data from cache", err)
		return DTO.AppStatus{}, models.NewError(500, "Server", "Internal server error")
	}

	appStatus, err := utils.UnmarshalData[DTO.AppStatus]([]byte(appStatusAsJson))
	if err != nil {
		as.LoggerService.Warn("Failed to unmarshal  data", err)
		return DTO.AppStatus{}, models.NewError(500, "Server", "Internal server error")
	}

	return *appStatus, nil
}

func (as *AppStatusService) checkAndCompareAppStatuses(ctx context.Context, cli *client.Client,
	appsToCheck []*models.AppToCheck,
) ([]DTO.AppStatus, []DTO.AppStatus) {
	appsStatusesChan := make(chan DTO.AppStatus, len(appsToCheck))
	appsToSendNotificationChan := make(chan DTO.AppStatus, len(appsToCheck))
	jobs := make(chan *models.AppToCheck, len(appsToCheck))

	workerCount := runtime.NumCPU()

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var appStatus DTO.AppStatus

				if job.IsDocker {
					container, err := cli.ContainerInspect(ctx, job.ID)
					if err != nil {
						as.LoggerService.Error("Failed to inspect container", err)
						continue
					}

					status := container.State.Status
					startedTime, err := time.Parse(time.RFC3339, container.State.StartedAt)
					if err != nil {
						as.LoggerService.Error("Failed to parse container start time", err)
						continue
					}

					duration := time.Since(startedTime)
					appStatus = *DTO.NewAppStatus(job.ID, status, startedTime, duration)
				} else {
					address := fmt.Sprintf("%s:%s", job.IpAddress, job.Port)
					conn, err := net.DialTimeout("tcp", address, 3*time.Second)
					status := "running"
					startedTime := time.Now()
					if err != nil {
						status = "stopped"
					}

					appStatus = *DTO.NewAppStatus(job.ID, status, startedTime, 0)
					if conn != nil {
						conn.Close()
					}
				}

				appsStatusesChan <- appStatus
				if appStatus.Status != job.Status {
					appsToSendNotificationChan <- appStatus
				}

				appStatusBytes, err := utils.MarshalData(appStatus)
				if err != nil {
					as.LoggerService.Error("Failed to marshal app status", map[string]any{"data": appStatus, "error": err.Error()})
					continue
				}

				if err := as.CacheService.SetData(ctx, "status-"+job.ID, string(appStatusBytes),
					2*time.Minute); err != nil {
					as.LoggerService.Error("Failed to set cache", map[string]any{"data": appStatus, "error": err.Error()})
				}
			}
		}()
	}

	for _, appToCheck := range appsToCheck {
		jobs <- appToCheck
	}
	close(jobs)

	wg.Wait()
	close(appsStatusesChan)
	close(appsToSendNotificationChan)

	var appsStatuses []DTO.AppStatus
	for appStatus := range appsStatusesChan {
		appsStatuses = append(appsStatuses, appStatus)
	}

	var appsToSendNotification []DTO.AppStatus
	for appToSendNotification := range appsToSendNotificationChan {
		appsToSendNotification = append(appsToSendNotification, appToSendNotification)
	}
	return appsStatuses, appsToSendNotification
}

func (as *AppStatusService) checkAppsStatus(ctx context.Context) ([]DTO.AppStatus, error) {
	appsToCheck, err := as.AppRepository.GetAppsToCheck(ctx)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.WithHost(as.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	appsStatuses, appsToSendNotification := as.checkAndCompareAppStatuses(ctx, cli, appsToCheck)
	if len(appsStatuses) > 0 {
		if err := as.AppRepository.InsertAppStatuses(ctx, appsStatuses); err != nil {
			as.LoggerService.Error("Failed to insert app statuses", err)
			return appsToSendNotification, err
		}
	}

	return appsToSendNotification, nil
}

func (as *AppStatusService) getAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error) {
	cacheKey := fmt.Sprintf("status-%s", id)

	doesAppStatusExists, err := as.CacheService.ExistsData(ctx, cacheKey)
	if err != nil {
		as.LoggerService.Warn("Failed to get info about data in cache", err)
	}

	if doesAppStatusExists > 0 {
		appStatus, err := as.readAppStatusFromCache(ctx, cacheKey)
		if err != nil {
			return DTO.AppStatus{}, err
		}
		return appStatus, nil
	}

	appStatus, err := as.AppRepository.GetAppStatus(ctx, id, ownerId)
	if err != nil {
		return DTO.AppStatus{}, err
	}

	return appStatus, nil
}
