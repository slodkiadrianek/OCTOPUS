package services

import (
	"context"
	"fmt"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AppNotificationsService struct {
	AppRepository appRepository
	LoggerService *utils.Logger
}

func NewAppNotificationsService(appRepository appRepository, loggerService *utils.Logger,
) *AppNotificationsService {
	return &AppNotificationsService{
		AppRepository: appRepository,
		LoggerService: loggerService,
	}
}

func (an *AppNotificationsService) assignNotificationToProperSendService(notificationsInfo []models.NotificationInfo) map[string][]models.
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
func (an *AppNotificationsService) sortNotificationsBySendToInformation(sortedNotificationsToSend map[string][]models.
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

func (an *AppNotificationsService) sendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	if len(appsStatuses) == 0 {
		return nil
	}

	an.LoggerService.Info("Started sending Notifications to users")

	notificationsInfo, err := an.AppRepository.GetUsersToSendNotifications(ctx, appsStatuses)
	if err != nil {
		return err
	}

	sortedNotificationsToSend := an.assignNotificationToProperSendService(notificationsInfo)
	discordNotifications, slackNotifications := an.sortNotificationsBySendToInformation(sortedNotificationsToSend)

	for webhookURL, message := range discordNotifications {
		payload := map[string]interface{}{
			"content":  message,
			"username": "OctopusBot",
		}

		body, err := utils.MarshalData(payload)
		if err != nil {
			return err
		}

		responseStatusCode, _, err := utils.DoHttpRequest(ctx, webhookURL, "", "POST", body, *an.LoggerService)
		if err != nil {
			return err
		}

		if responseStatusCode >= 300 {
			an.LoggerService.Warn("Webhook returned non-success status", "status", responseStatusCode)
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

		responseStatusCode, _, err := utils.DoHttpRequest(ctx, webhookURL, "", "POST", body, *an.LoggerService)
		if err != nil {
			return err
		}

		if responseStatusCode >= 300 {
			an.LoggerService.Warn("Webhook returned non-success status", "status", responseStatusCode)
		}
	}
	return nil
}
