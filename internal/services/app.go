package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

//func (a *AppService) doHttpRequest(ctx context.Context, url, authorizationHeader, method string, body []byte) (int,
//	map[string]any, error) {
//	httpClient := &http.Client{}
//	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
//	if err != nil {
//		a.LoggerService.Error("Failed to create webhook request", err)
//		return 0, map[string]any{}, err
//	}
//	req.Header.Add("Authorization", authorizationHeader)
//
//	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//	response, err := httpClient.Do(req)
//	if err != nil {
//		a.LoggerService.Error("Failed to send webhook request", err)
//		return 0, map[string]any{}, err
//	}
//	defer response.Body.Close()
//	var bodyFromResponse map[string]any
//	err = json.NewDecoder(response.Body).Decode(&bodyFromResponse)
//	return response.StatusCode, bodyFromResponse, nil
//}

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

func (a *AppService) DeleteApp(ctx context.Context, id string, ownerId int) error {
	err := a.AppRepository.DeleteApp(ctx, id, ownerId)
	if err != nil {
		return err
	}
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
	map[string]string,
) {
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
		responseStatusCode, _, err := as.doHttpRequest(ctx, webhookURL, "", "POST", body)
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
