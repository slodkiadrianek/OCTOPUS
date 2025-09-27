package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AppService struct {
	AppRepository *repository.AppRepository
	Logger        *utils.Logger
	CacheService  *config.CacheService
	DockerHost    string
}

func NewAppService(appRepository *repository.AppRepository, logger *utils.Logger, cacheService *config.CacheService, dockerHost string) *AppService {
	return &AppService{
		AppRepository: appRepository,
		Logger:        logger,
		CacheService:  cacheService,
		DockerHost:    dockerHost,
	}
}

func (a *AppService) CreateApp(ctx context.Context, app DTO.CreateApp, ownerId int) error {
	id, err := utils.GenerateID()
	if err != nil {
		return err
	}
	appDto := DTO.NewApp(id, app.Name, app.Description, false, ownerId, app.IpAddress, app.Port)
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

func (a *AppService) GetAppStatus(ctx context.Context, id string) (DTO.AppStatus, error) {
	cacheKey := fmt.Sprintf("status-%s", id)
	doesExist, err := a.CacheService.ExistsData(ctx, cacheKey)
	if err != nil {
		a.Logger.Warn("Failed to get info about data in cache", err)
	}
	if doesExist > 0 {
		data, err := a.CacheService.GetData(ctx, cacheKey)
		if err != nil {
			a.Logger.Warn("Failed to get data from cache", err)
			return DTO.AppStatus{}, models.NewError(500, "Server", "Internal server error")
		}
		appStatus, err := utils.UnmarshalData[DTO.AppStatus]([]byte(data))
		if err != nil {
			a.Logger.Warn("Failed to unmarshal  data", err)
			return DTO.AppStatus{}, models.NewError(500, "Server", "Internal server error")
		}
		return *appStatus, nil
	}
	appStatus, err := a.AppRepository.GetAppStatus(ctx, id)
	if err != nil {
		return DTO.AppStatus{}, err
	}
	return appStatus, nil
}

func (a *AppService) DeleteApp(ctx context.Context, id string) error {
	err := a.AppRepository.DeleteApp(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetLogs(ctx context.Context, appId string) (string, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(a.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: true,
		Since:      "",
		Until:      "",
		Tail:       "100",
	}
	reader, err := cli.ContainerLogs(ctx, appId, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return "", err
	}
	logs := utils.StripANSI(buf.String())
	return logs, nil
}

func (a *AppService) CheckAppsStatus(ctx context.Context) ([]DTO.AppStatus, error) {
	apps, err := a.AppRepository.GetAppsToCheck(ctx)
	if err != nil {
		return []DTO.AppStatus{}, err
	}
	workerCount := runtime.NumCPU()
	cli, err := client.NewClientWithOpts(client.WithHost(a.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return []DTO.AppStatus{}, err
	}
	var appsStatuses []DTO.AppStatus
	var appsToSendNotification []DTO.AppStatus
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
					if status != job.Status {
						appsToSendNotification = append(appsToSendNotification, appStatus)
					}
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
					address := fmt.Sprintf("%s:%s", job.IpAddress, job.Port)
					conn, err := net.DialTimeout("tcp", address, time.Second*3)
					status := "running"
					duration := time.Since(time.Now())
					changedAt := time.Now()
					if err != nil {
						a.Logger.Error("Failed to check status inside of a container", err)
						status = "stopped"
						appStatus := *DTO.NewAppStatus(job.Id, status, changedAt, duration)
						appsStatuses = append(appsStatuses, appStatus)
						if status != job.Status {
							appsToSendNotification = append(appsToSendNotification, appStatus)
						}
						continue
					}
					appStatus := *DTO.NewAppStatus(job.Id, status, changedAt, duration)
					bodyBytes, err := utils.MarshalData(appStatus)
					if err != nil {
						a.Logger.Error("Failed to convert data to json", map[string]any{
							"data":  appStatus,
							"error": err.Error(),
						})
						continue
					}
					if status != job.Status {
						appsToSendNotification = append(appsToSendNotification, appStatus)
					}
					appsStatuses = append(appsStatuses, appStatus)
					defer conn.Close()
					err = a.CacheService.SetData(ctx, "status-"+job.Id, string(bodyBytes), time.Minute*2)
					if err != nil {
						a.Logger.Error("Failed to set data in cache", map[string]any{
							"data":  appStatus,
							"error": err.Error(),
						})
						continue
					}
				}
			}
		}(i + 1)
	}
	for _, app := range apps {
		jobs <- app
	}
	close(jobs)
	wg.Wait()
	if len(appsStatuses) > 0 {
		err = a.AppRepository.InsertAppStatuses(ctx, appsStatuses)
	}
	return appsToSendNotification, err
}

func (a *AppService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	if len(appsStatuses) == 0 {
		return nil
	}
	a.Logger.Info("Started sending Notfications to user")
	appsToSendNotifications, err := a.AppRepository.GetUsersToSendNotifications(ctx, appsStatuses)
	if err != nil {
		return err
	}
	sortedData := map[string][]models.SendNotificationTo{
		"Discord": {},
		"Slack":   {},
		"Email":   {},
	}
	for _, app := range appsToSendNotifications {
		if app.DiscordNotifications {
			sortedData["Discord"] = append(sortedData["Discord"], app)
		}
		if app.SlackNotifications {
			sortedData["Slack"] = append(sortedData["Slack"], app)
		}
		if app.EmailNotifications {
			sortedData["Email"] = append(sortedData["Email"], app)
		}
	}
	sortedDataDiscord := map[string]string{}
	sortedDataSlack := map[string]string{}
	// sortedDataEmail := map[string]string{}
	for _, val := range sortedData["Discord"] {
		sortedDataDiscord[val.DiscordWebhook] += fmt.Sprintf("%s-%s-%s\n", val.Id, val.Name, val.Status)
	}
	for _, val := range sortedData["Slack"] {
		sortedDataSlack[val.SlackWebhook] += fmt.Sprintf("%s-%s-%s\n", val.Id, val.Name, val.Status)
	}
	// for _, val := range sortedData["Email"] {
	// 	sortedDataDiscord[val.Email] += fmt.Sprintf("%s-%s-%s", val.Id, val.Name, val.Status)
	// }
	for i, val := range sortedDataDiscord {
		payload := map[string]interface{}{
			"content":  val,
			"username": "OctopusBot",
			"embeds":   []interface{}{},
		}

		jsonData, err := utils.MarshalData(payload)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", i, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client := &http.Client{}
		client.Do(req)
	}
	for i, val := range sortedDataSlack {
		payload := map[string]interface{}{
			"text":     val,
			"username": "OctopusBot",
		}

		jsonData, err := utils.MarshalData(payload)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", i, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		client.Do(req)
	}
	return nil
}

func (a *AppService) UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp) error {
	err := a.AppRepository.UpdateApp(ctx, appId, app)
	if err != nil {
		return err
	}
	return nil
}
