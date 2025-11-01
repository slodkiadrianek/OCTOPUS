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
	GetUsersToSendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) ([]models.NotificationInfo, error)
}
type AppService struct {
	AppRepository   appRepository
	LoggerService   *utils.Logger
	CacheService    CacheService
	DockerHost      string
	RouteRepository *repository.RouteRepository
}

func NewAppService(appRepository appRepository, loggerService *utils.Logger, cacheService CacheService,
	dockerHost string, routeRepository *repository.RouteRepository,
) *AppService {
	return &AppService{
		AppRepository:   appRepository,
		LoggerService:   loggerService,
		CacheService:    cacheService,
		DockerHost:      dockerHost,
		RouteRepository: routeRepository,
	}
}
func (a *AppService) doHttpRequest(ctx context.Context, url, authorizationHeader, method string, body []byte) (int,
	map[string]any, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		a.LoggerService.Error("Failed to create webhook request", err)
		return 0, map[string]any{}, err
	}
	req.Header.Add("Authorization", authorizationHeader)

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := httpClient.Do(req)
	if err != nil {
		a.LoggerService.Error("Failed to send webhook request", err)
		return 0, map[string]any{}, err
	}
	defer response.Body.Close()
	var bodyFromResponse map[string]any
	err = json.NewDecoder(response.Body).Decode(&bodyFromResponse)
	return response.StatusCode, bodyFromResponse, nil
}

func (a *AppService) CreateApp(ctx context.Context, app DTO.CreateApp, ownerId int) error {
	GeneratedId, err := utils.GenerateID()
	if err != nil {
		return err
	}
	appDto := DTO.NewApp(GeneratedId, app.Name, app.Description, false, ownerId, app.IpAddress, app.Port)
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

func (a *AppService) readAppStatusFromCache(ctx context.Context, cacheKey string) (DTO.AppStatus, error) {
	appStatusAsJson, err := a.CacheService.GetData(ctx, cacheKey)
	if err != nil {
		a.LoggerService.Warn("Failed to get data from cache", err)
		return DTO.AppStatus{}, models.NewError(500, "Server", "Internal server error")
	}
	appStatus, err := utils.UnmarshalData[DTO.AppStatus]([]byte(appStatusAsJson))
	if err != nil {
		a.LoggerService.Warn("Failed to unmarshal  data", err)
		return DTO.AppStatus{}, models.NewError(500, "Server", "Internal server error")
	}
	return *appStatus, nil
}

func (a *AppService) GetAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error) {
	cacheKey := fmt.Sprintf("status-%s", id)
	doesAppStatusExists, err := a.CacheService.ExistsData(ctx, cacheKey)
	if err != nil {
		a.LoggerService.Warn("Failed to get info about data in cache", err)
	}
	if doesAppStatusExists > 0 {
		appStatus, err := a.readAppStatusFromCache(ctx, cacheKey)
		if err != nil {
			return DTO.AppStatus{}, err
		}
		return appStatus, nil
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

func (a *AppService) checkAndCompareAppStatuses(ctx context.Context, cli *client.Client,
	appsToCheck []*models.AppToCheck) ([]DTO.AppStatus, []DTO.AppStatus) {
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
						a.LoggerService.Error("Failed to inspect container", err)
						continue
					}
					status := container.State.Status
					startedTime, err := time.Parse(time.RFC3339, container.State.StartedAt)
					if err != nil {
						a.LoggerService.Error("Failed to parse container start time", err)
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
					a.LoggerService.Error("Failed to marshal app status", map[string]any{"data": appStatus, "error": err.Error()})
					continue
				}
				if err := a.CacheService.SetData(ctx, "status-"+job.ID, string(appStatusBytes),
					2*time.Minute); err != nil {
					a.LoggerService.Error("Failed to set cache", map[string]any{"data": appStatus, "error": err.Error()})
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

func (a *AppService) CheckAppsStatus(ctx context.Context) ([]DTO.AppStatus, error) {
	appsToCheck, err := a.AppRepository.GetAppsToCheck(ctx)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.WithHost(a.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	appsStatuses, appsToSendNotification := a.checkAndCompareAppStatuses(ctx, cli, appsToCheck)
	if len(appsStatuses) > 0 {
		if err := a.AppRepository.InsertAppStatuses(ctx, appsStatuses); err != nil {
			a.LoggerService.Error("Failed to insert app statuses", err)
			return appsToSendNotification, err
		}
	}
	return appsToSendNotification, nil
}

func (a *AppService) sortRoutesToTest(routesToTest []models.RouteToTest) map[string][]models.RouteToTest {
	sortedRoutesToTests := make(map[string][]models.RouteToTest)
	for _, routeToTest := range routesToTest {
		key := routeToTest.Name + routeToTest.AppId
		if routeToTest.ParentID == 0 {
			sortedRoutesToTests[key] = append([]models.RouteToTest{routeToTest},
				sortedRoutesToTests[key]...)
		} else {
			sortedRoutesToTests[key] = append(sortedRoutesToTests[key], routeToTest)
		}
	}
	return sortedRoutesToTests
}

func (a *AppService) prepareRouteDataForTestRequest(route models.RouteToTest) (string, string, []byte, error) {
	authorizationHeader := "Bearer " + route.RequestAuthorization
	splittedPath := strings.Split(route.Path, "/")
	for i := 0; i < len(splittedPath); i++ {
		partOfPath := splittedPath[i]
		leftBrace := strings.Contains(partOfPath, "{")
		rightBrace := strings.Contains(partOfPath, "}")
		if leftBrace && rightBrace {
			param := partOfPath[1 : len(partOfPath)-1]
			splittedPath[i] = route.RequestParams[param]
		}
	}
	var query []string
	for key, val := range route.RequestQuery {
		query = append(query, key+"="+val)
	}

	path := strings.Join(splittedPath, "/")
	url := "http://" + route.IpAddress + ":" + route.Port + path + "?" + strings.Join(query, "&")
	jsonData, err := utils.MarshalData(route.RequestBody)
	if err != nil {
		a.LoggerService.Error("Failed to marshal webhook payload", err)
		return "", "", []byte{}, err
	}
	return authorizationHeader, url, jsonData, nil
}

func (a *AppService) prepareDataForTheNextRoute(route models.RouteToTest, key string,
	val any) (map[string]any, map[string]string, map[string]string, string, string) {
	routeStatus := "unknown"
	nextRouteBody := make(map[string]any)
	nextRouteParams := make(map[string]string)
	nextRouteQuery := make(map[string]string)
	nextRouteAuthorizationHeader := ""

	if slices.Contains(route.NextRouteBody, key) {
		nextRouteBody[key] = val
	}
	if slices.Contains(route.NextRouteParams, key) {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for param"
			return map[string]any{}, map[string]string{}, map[string]string{}, "", routeStatus
		}
		nextRouteParams[key] = valueConvertedToString
	}
	if slices.Contains(route.NextRouteQuery, key) {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for query"
			return map[string]any{}, map[string]string{}, map[string]string{}, "", routeStatus
		}
		nextRouteQuery[key] = valueConvertedToString
	}
	valueConvertedToString, ok := val.(string)
	if !ok {
		routeStatus = "Failed;Wrong type of the property for authorization header"
		return map[string]any{}, map[string]string{}, map[string]string{}, "", routeStatus
	}
	if strings.Contains(valueConvertedToString, "eyJlbWFpbCI6IlRFU1QiLCJleHAiOjE3N") {
		nextRouteAuthorizationHeader = valueConvertedToString
	}
	return nextRouteBody, nextRouteParams, nextRouteQuery, nextRouteAuthorizationHeader, routeStatus
}

func (a *AppService) CheckRoutesStatus(ctx context.Context) error {
	a.LoggerService.Info("Started checking statuses of the routes")
	routesToTest, err := a.RouteRepository.GetWorkingRoutesToTest(ctx)
	if err != nil {
		return err
	}
	if len(routesToTest) < 1 {
		return nil
	}
	sortedRoutesToTests := a.sortRoutesToTest(routesToTest)
	routesStatuses := make(map[int]string)
	for _, routesToTest := range sortedRoutesToTests {
		nextRouteBody := make(map[string]any)
		nextRouteParams := make(map[string]string)
		nextRouteQuery := make(map[string]string)
		nextRouteAuthorizationHeader := ""
		for _, route := range routesToTest {
			routeStatus := "unknown"
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
			authorizationHeader, url, body, err := a.prepareRouteDataForTestRequest(route)
			if err != nil {
				return err
			}
			responseStatusCode, responseBody, err := a.doHttpRequest(ctx, url, authorizationHeader, route.Method, body)
			if len(responseBody) != len(route.ResponseBody) {
				routeStatus = "Failed;Different body"
				routesStatuses[route.ID] = routeStatus
				break
			}
			if responseStatusCode != route.ResponseStatusCode {
				routeStatus = "Failed;Status Code"
				routesStatuses[route.ID] = routeStatus
				break
			}
			for key, val := range responseBody {
				nextRouteBody, nextRouteParams, nextRouteQuery, nextRouteAuthorizationHeader,
					routeStatus = a.prepareDataForTheNextRoute(route, key, val)
			}
			routeStatus = "success"
			routesStatuses[route.ID] = routeStatus
		}

	}
	a.LoggerService.Info("The routes statuses have started inserting into database", routesStatuses)
	err = a.RouteRepository.UpdateWorkingRoutesStatuses(ctx, routesStatuses)
	if err != nil {
		return err
	}
	a.LoggerService.Info("The route statuses have finished inserting into the database.", routesStatuses)
	return nil
}

func (a *AppService) assignNotificationToProperSendService(notificationsInfo []models.NotificationInfo) map[string][]models.
	NotificationInfo {
	sortedNotificationsToSend := map[string][]models.NotificationInfo{
		"Discord": {},
		"Slack":   {},
		"Email":   {},
	}

	for _, notificationInfo := range notificationsInfo {
		if notificationInfo.DiscordNotificationsSettings {
			sortedNotificationsToSend["Discord"] = append(sortedNotificationsToSend["Discord"], notificationInfo)
		}
		if notificationInfo.SlackNotificationsSettings {
			sortedNotificationsToSend["Slack"] = append(sortedNotificationsToSend["Slack"], notificationInfo)
		}
		if notificationInfo.EmailNotificationsSettings {
			sortedNotificationsToSend["Email"] = append(sortedNotificationsToSend["Email"], notificationInfo)
		}
	}
	return sortedNotificationsToSend
}

// Send to information is, for example, email, discord webhook url or slack webhook url
func (a *AppService) sortNotificationsBySendToInformation(sortedNotificationsToSend map[string][]models.
	NotificationInfo) (
	map[string]string,
	map[string]string) {
	discordNotifications := map[string]string{}
	slackNotifications := map[string]string{}

	for _, discordNotificationInfo := range sortedNotificationsToSend["Discord"] {
		discordNotifications[discordNotificationInfo.DiscordWebhookUrl] += fmt.Sprintf("%s - %s - %s\n",
			discordNotificationInfo.ID, discordNotificationInfo.Name, discordNotificationInfo.Status)
	}

	for _, slackNotificationInfo := range sortedNotificationsToSend["Slack"] {
		slackNotifications[slackNotificationInfo.SlackWebhookUrl] += fmt.Sprintf("%s - %s - %s\n",
			slackNotificationInfo.ID, slackNotificationInfo.Name, slackNotificationInfo.Status)
	}
	return discordNotifications, slackNotifications
}

func (a *AppService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	if len(appsStatuses) == 0 {
		return nil
	}
	a.LoggerService.Info("Started sending Notifications to users")
	notificationsInfo, err := a.AppRepository.GetUsersToSendNotifications(ctx, appsStatuses)
	if err != nil {
		return err
	}
	sortedNotificationsToSend := a.assignNotificationToProperSendService(notificationsInfo)
	discordNotifications, slackNotifications := a.sortNotificationsBySendToInformation(sortedNotificationsToSend)

	for webhookURL, message := range discordNotifications {
		payload := map[string]interface{}{
			"content":  message,
			"username": "OctopusBot",
		}
		body, err := utils.MarshalData(payload)
		if err != nil {
			return err
		}
		responseStatusCode, _, err := a.doHttpRequest(ctx, webhookURL, "", "POST", body)
		if err != nil {
			return err
		}
		if responseStatusCode >= 300 {
			a.LoggerService.Warn("Webhook returned non-success status", "status", responseStatusCode)
		}
	}

	for webhookURL, message := range slackNotifications {
		payload := map[string]interface{}{
			"content":  message,
			"username": "OctopusBot",
		}
		body, err := utils.MarshalData(payload)
		if err != nil {
			return err
		}
		responseStatusCode, _, err := a.doHttpRequest(ctx, webhookURL, "", "POST", body)
		if err != nil {
			return err
		}
		if responseStatusCode >= 300 {
			a.LoggerService.Warn("Webhook returned non-success status", "status", responseStatusCode)
		}
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
