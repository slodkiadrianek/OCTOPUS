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

func (a *AppService) GetApp(ctx context.Context, id string) (*models.App, error) {
	app, err := a.AppRepository.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) GetApps(ctx context.Context) ([]models.App, error) {
	apps, err := a.AppRepository.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	return apps, nil
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
		return nil, err
	}

	workerCount := runtime.NumCPU()
	cli, err := client.NewClientWithOpts(client.WithHost(a.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	appsStatusesChan := make(chan DTO.AppStatus, len(apps))
	appsToSendChan := make(chan DTO.AppStatus, len(apps))

	jobs := make(chan *models.AppToCheck, len(apps))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var appStatus DTO.AppStatus

				if job.IsDocker {
					container, err := cli.ContainerInspect(ctx, job.Id)
					if err != nil {
						a.Logger.Error("Failed to inspect container", err)
						continue
					}

					status := container.State.Status
					startedTime, err := time.Parse(time.RFC3339, container.State.StartedAt)
					if err != nil {
						a.Logger.Error("Failed to parse container start time", err)
						continue
					}

					duration := time.Since(startedTime)
					appStatus = *DTO.NewAppStatus(job.Id, status, startedTime, duration)
				} else {
					address := fmt.Sprintf("%s:%s", job.IpAddress, job.Port)
					conn, err := net.DialTimeout("tcp", address, 3*time.Second)
					status := "running"
					startedTime := time.Now()
					if err != nil {
						status = "stopped"
					}
					appStatus = *DTO.NewAppStatus(job.Id, status, startedTime, 0)
					if conn != nil {
						conn.Close()
					}
				}

				appsStatusesChan <- appStatus
				fmt.Println(appStatus.Status, job.Status)
				if appStatus.Status != job.Status {
					appsToSendChan <- appStatus
				}

				bodyBytes, err := utils.MarshalData(appStatus)
				if err != nil {
					a.Logger.Error("Failed to marshal app status", map[string]any{"data": appStatus, "error": err.Error()})
					continue
				}
				if err := a.CacheService.SetData(ctx, "status-"+job.Id, string(bodyBytes), 2*time.Minute); err != nil {
					a.Logger.Error("Failed to set cache", map[string]any{"data": appStatus, "error": err.Error()})
				}
			}
		}()
	}

	for _, app := range apps {
		jobs <- app
	}
	close(jobs)

	wg.Wait()
	close(appsStatusesChan)
	close(appsToSendChan)

	var appsStatuses []DTO.AppStatus
	for status := range appsStatusesChan {
		appsStatuses = append(appsStatuses, status)
	}

	var appsToSendNotification []DTO.AppStatus
	for notify := range appsToSendChan {
		appsToSendNotification = append(appsToSendNotification, notify)
	}

	if len(appsStatuses) > 0 {
		if err := a.AppRepository.InsertAppStatuses(ctx, appsStatuses); err != nil {
			a.Logger.Error("Failed to insert app statuses", err)
			return appsToSendNotification, err
		}
	}

	return appsToSendNotification, nil
}

func (a *AppService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	fmt.Println("IS STATUS", appsStatuses)
	if len(appsStatuses) == 0 {
		return nil
	}

	a.Logger.Info("Started sending Notifications to users")

	// Fetch users to notify
	appsToSendNotifications, err := a.AppRepository.GetUsersToSendNotifications(ctx, appsStatuses)
	if err != nil {
		return err
	}

	// Group apps by notification type
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

	// Aggregate messages per webhook
	discordMessages := map[string]string{}
	slackMessages := map[string]string{}

	for _, val := range sortedData["Discord"] {
		discordMessages[val.DiscordWebhook] += fmt.Sprintf("%s - %s - %s\n", val.Id, val.Name, val.Status)
	}

	for _, val := range sortedData["Slack"] {
		slackMessages[val.SlackWebhook] += fmt.Sprintf("%s - %s - %s\n", val.Id, val.Name, val.Status)
	}

	// Reuse HTTP client
	client := &http.Client{}

	// Helper function to send POST requests with context
	sendWebhook := func(ctx context.Context, url string, payload map[string]interface{}) {
		jsonData, err := utils.MarshalData(payload)
		if err != nil {
			a.Logger.Error("Failed to marshal webhook payload", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			a.Logger.Error("Failed to create webhook request", err)
			return
		}

		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Do(req)
		if err != nil {
			a.Logger.Error("Failed to send webhook request", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			a.Logger.Warn("Webhook returned non-success status", "status", resp.Status)
		}
	}

	// Send Discord notifications
	for webhookURL, message := range discordMessages {
		payload := map[string]interface{}{
			"content":  message,
			"username": "OctopusBot",
		}
		sendWebhook(ctx, webhookURL, payload)
	}

	// Send Slack notifications
	for webhookURL, message := range slackMessages {
		payload := map[string]interface{}{
			"text": message,
		}
		sendWebhook(ctx, webhookURL, payload)
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
