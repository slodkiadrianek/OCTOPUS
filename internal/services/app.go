package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type appRepository interface {
	InsertApp(ctx context.Context, app []DTO.App) error
	GetApp(ctx context.Context, id string, ownerId int) (*models.App, error)
	GetApps(ctx context.Context, ownerId int) ([]models.App, error)
	DeleteApp(ctx context.Context, id string, ownerId int) error
	GetAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error)
	GetAppsToCheck(ctx context.Context) ([]*models.AppToCheck, error)
	UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error
	InsertAppStatuses(ctx context.Context, appsStatuses []DTO.AppStatus) error
	GetUsersToSendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) ([]models.SendNotificationTo, error)
}
type AppService struct {
	AppRepository   appRepository
	Logger          *utils.Logger
	CacheService    CacheService
	DockerHost      string
	RouteRepository *repository.RouteRepository
}

func NewAppService(appRepository appRepository, logger *utils.Logger, cacheService CacheService,
	dockerHost string, routeRepository *repository.RouteRepository,
) *AppService {
	return &AppService{
		AppRepository:   appRepository,
		Logger:          logger,
		CacheService:    cacheService,
		DockerHost:      dockerHost,
		RouteRepository: routeRepository,
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

func (a *AppService) GetApp(ctx context.Context, id string, ownerId int) (*models.App, error) {
	app, err := a.AppRepository.GetApp(ctx, id, ownerId)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) GetApps(ctx context.Context, ownerId int) ([]models.App, error) {
	apps, err := a.AppRepository.GetApps(ctx, ownerId)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (a *AppService) GetAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error) {
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
	appStatus, err := a.AppRepository.GetAppStatus(ctx, id, ownerId)
	if err != nil {
		return DTO.AppStatus{}, err
	}
	return appStatus, nil
}

func (a *AppService) DeleteApp(ctx context.Context, id string, ownerId int) error {
	err := a.AppRepository.DeleteApp(ctx, id, ownerId)
	if err != nil {
		return err
	}
	return nil
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

func (a *AppService) CheckRoutesStatus(ctx context.Context) error {
	routesToTest, err := a.RouteRepository.GetWorkingRoutesToTest(ctx)
	if err != nil {
		return err
	}
	var sortedRoutesToTests map[string][]DTO.RouteToTest
	for _, routeToTest := range routesToTest {
		key := routeToTest.Name + routeToTest.AppId
		if routeToTest.ParentId == 0 {
			sortedRoutesToTests[key] = append([]DTO.RouteToTest{routeToTest},
				sortedRoutesToTests[key]...)
		} else {
			sortedRoutesToTests[key] = append([]DTO.RouteToTest{routeToTest},
				sortedRoutesToTests[key]...)
		}
	}
	var routesStatuses map[int]string
	client := &http.Client{}
	for _, routesToTest := range sortedRoutesToTests {
		var nextRouteBody map[string]any
		var nextRouteParams map[string]string
		var nextRouteQuery map[string]string
		var nextRouteAuthorizationHeader string
		for _, route := range routesToTest {
			var routeStatus string
			if len(nextRouteBody) > 0 {
				route.RequestBody = nextRouteBody
			}
			if len(nextRouteParams) > 0 {
				route.RequestParams = nextRouteParams
			}
			if len(nextRouteQuery) > 0 {
				route.RequestQuery = nextRouteQuery
			}
			if len(nextRouteAuthorizationHeader) > 0 {
				route.RequestAuthorization = nextRouteAuthorizationHeader
			}
			authorizationHeader := "Bearer " + route.RequestAuthorization
			splittedPath := strings.Split(route.Path, "/")
			for _, val := range splittedPath {
				leftBrace := strings.Contains(val, "{")
				rightBrace := strings.Contains(val, "}")
				if leftBrace && rightBrace {
					param := val[1 : len(val)-1]
					val = route.RequestParams[param]
				}
			}
			var query []string
			for key, val := range route.RequestQuery {
				query = append(query, key+"="+val)
			}

			path := strings.Join(splittedPath, "/")
			url := route.IpAddress + ":" + route.Port + path + "?" + strings.Join(query, "&")
			jsonData, err := utils.MarshalData(route.RequestBody)
			if err != nil {
				a.Logger.Error("Failed to marshal webhook payload", err)
				return err
			}

			req, err := http.NewRequestWithContext(ctx, route.Method, url, bytes.NewBuffer(jsonData))
			req.Header.Add("Authorization", authorizationHeader)
			if err != nil {
				a.Logger.Error("Failed to create webhook request", err)
				return err
			}

			req.Header.Set("Content-Type", "application/json; charset=UTF-8")

			resp, err := client.Do(req)
			if err != nil {
				a.Logger.Error("Failed to send webhook request", err)
				return err
			}
			defer resp.Body.Close()
			var body map[string]any
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err != nil {
				a.Logger.Error("Failed to read body from the request", err)
				return err
			}
			if len(body) != len(route.ResponseBody) {
				routeStatus = "Failed;Different body"
				routesStatuses[route.Id] = routeStatus
				break
			}
			if resp.StatusCode != route.ResponseStatusCode {
				routeStatus = "Failed;Status Code"
				routesStatuses[route.Id] = routeStatus
				break
			}
			for key, val := range body {
				if slices.Contains(route.NextRouteBody, key) {
					nextRouteBody[key] = val
				}
				if slices.Contains(route.NextRouteParams, key) {
					valS, ok := val.(string)
					if !ok {
						routeStatus = "Failed;Wrong type of the property for param"
						routesStatuses[route.Id] = routeStatus
						break
					}
					nextRouteParams[key] = valS
				}
				if slices.Contains(route.NextRouteQuery, key) {
					valS, ok := val.(string)
					if !ok {
						routeStatus = "Failed;Wrong type of the property for query"
						routesStatuses[route.Id] = routeStatus
						break
					}
					nextRouteQuery[key] = valS
				}
				if strings.Contains(route.NextAuthorizationHeader, key) {

					valS, ok := val.(string)
					if !ok {
						routeStatus = "Failed;Wrong type of the property for authorization header"
						routesStatuses[route.Id] = routeStatus
						break
					}
					nextRouteAuthorizationHeader = valS
				}
			}
			routeStatus = "success"
			routesStatuses[route.Id] = routeStatus
		}

	}
	fmt.Println(routesStatuses)
	return nil
}

func (a *AppService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	if len(appsStatuses) == 0 {
		return nil
	}
	a.Logger.Info("Started sending Notifications to users")
	usersToSendNotifications, err := a.AppRepository.GetUsersToSendNotifications(ctx, appsStatuses)
	if err != nil {
		return err
	}

	sortedData := map[string][]models.SendNotificationTo{
		"Discord": {},
		"Slack":   {},
		"Email":   {},
	}

	for _, app := range usersToSendNotifications {
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

	discordMessages := map[string]string{}
	slackMessages := map[string]string{}

	for _, val := range sortedData["Discord"] {
		discordMessages[val.DiscordWebhook] += fmt.Sprintf("%s - %s - %s\n", val.Id, val.Name, val.Status)
	}

	for _, val := range sortedData["Slack"] {
		slackMessages[val.SlackWebhook] += fmt.Sprintf("%s - %s - %s\n", val.Id, val.Name, val.Status)
	}

	client := &http.Client{}

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

	for webhookURL, message := range discordMessages {
		payload := map[string]interface{}{
			"content":  message,
			"username": "OctopusBot",
		}
		sendWebhook(ctx, webhookURL, payload)
	}

	for webhookURL, message := range slackMessages {
		payload := map[string]interface{}{
			"text": message,
		}
		sendWebhook(ctx, webhookURL, payload)
	}

	return nil
}

func (a *AppService) UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error {
	err := a.AppRepository.UpdateApp(ctx, appId, app, ownerId)
	if err != nil {
		return err
	}
	return nil
}
