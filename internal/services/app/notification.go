package servicesApp

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
)

type AppNotificationsService struct {
	appRepository interfaces.AppRepository
	loggerService utils.LoggerService
}

func NewAppNotificationsService(appRepository interfaces.AppRepository, loggerService utils.LoggerService,
) *AppNotificationsService {
	return &AppNotificationsService{
		appRepository: appRepository,
		loggerService: loggerService,
	}
}

func (an *AppNotificationsService) assignNotificationToProperSendService(notificationsInfo []models.NotificationInfo) map[string][]models.
	NotificationInfo {

	sortedNotificationsToSend := make(map[string][]models.NotificationInfo, len(notificationsInfo))

	for _, notificationInfo := range notificationsInfo {
		if notificationInfo.DiscordNotificationsSettings && notificationInfo.DiscordWebhookUrl != "" {
			sortedNotificationsToSend["Discord"] = append(sortedNotificationsToSend["Discord"], notificationInfo)
		}
		if notificationInfo.SlackNotificationsSettings && notificationInfo.SlackWebhookUrl != "" {
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
	discordNotifications := make(map[string]string, len(sortedNotificationsToSend))
	slackNotifications := make(map[string]string, len(sortedNotificationsToSend))

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

func (an *AppNotificationsService) sendWebhook(ctx context.Context, notifications map[string]string) error {
	var wg sync.WaitGroup
	notificationsChan := make(chan map[string]string, len(notifications))
	type WebhookJob struct {
		url     string
		message string
	}
	jobs := make(chan WebhookJob, len(notifications))
	workerCount := runtime.NumCPU()
	errorChan := make(chan error)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var payload map[string]any
				if strings.Contains(job.url, "slack") {
					payload = map[string]any{
						"text":     job.message,
						"username": "OctopusBot",
					}
				} else {
					payload = map[string]any{
						"content":  job.message,
						"username": "OctopusBot",
					}
				}

				body, err := utils.MarshalData(payload)
				if err != nil {
					an.loggerService.Error("Failed to unmarshal a data", err)
					errorChan <- err
					continue
				}

				responseStatusCode, _, err := request.SendHttp(ctx, job.url, "", "POST", body, false)
				if err != nil {
					an.loggerService.Info("Failed to send a webhook", err)
					errorChan <- err
					continue
				}

				if responseStatusCode >= 300 {
					an.loggerService.Info("Webhook returned non-success status", "status", responseStatusCode)
					continue
				}
			}
		}()
	}
	for webhookURL, message := range notifications {
		jobs <- WebhookJob{url: webhookURL, message: message}
	}

	close(jobs)
	wg.Wait()
	close(notificationsChan)
	close(errorChan)

	select {
	case err := <-errorChan:
		return err
	default:
		return nil
	}
}

func (an *AppNotificationsService) SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	if len(appsStatuses) == 0 {
		return nil
	}

	an.loggerService.Info("Started sending notifications to users")

	notificationsInfo, err := an.appRepository.GetUsersToSendNotifications(ctx, appsStatuses)
	if err != nil {
		return err
	}

	sortedNotificationsToSend := an.assignNotificationToProperSendService(notificationsInfo)
	discordNotifications, slackNotifications := an.sortNotificationsBySendToInformation(sortedNotificationsToSend)
	var discordWebhookError, slackWebhookError error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		discordWebhookError = an.sendWebhook(ctx, discordNotifications)
	}()
	go func() {
		defer wg.Done()
		slackWebhookError = an.sendWebhook(ctx, slackNotifications)
	}()

	wg.Wait()
	if discordWebhookError != nil {
		return discordWebhookError
	}
	if slackWebhookError != nil {
		return slackWebhookError
	}
	an.loggerService.Info("Successfully sent notifications to user")

	return nil
}
