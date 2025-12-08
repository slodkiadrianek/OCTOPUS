package servicesApp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppNotificationsService_assignNotificationToProperSendService(t *testing.T) {
	type args struct {
		name                        string
		notifications               []models.NotificationInfo
		setupMock                   func() interfaces.AppRepository
		expectedSortedNotifications map[string][]models.NotificationInfo
	}
	testsScenarios := []args{
		{
			name: "assign notification to proper sender service",
			notifications: []models.NotificationInfo{
				{
					DiscordNotificationsSettings: true,
					DiscordWebhookUrl:            "https://webhook.example.com",
				},
				{
					SlackNotificationsSettings: true,
					SlackWebhookUrl:            "https://webhook.example.com",
				},
				{
					EmailNotificationsSettings: true,
					Email:                      "test@gmail.com",
				},
			},
			expectedSortedNotifications: map[string][]models.NotificationInfo{
				"Discord": []models.NotificationInfo{
					{
						DiscordNotificationsSettings: true,
						DiscordWebhookUrl:            "https://webhook.example.com",
					},
				},
				"Slack": []models.NotificationInfo{

					{
						SlackNotificationsSettings: true,
						SlackWebhookUrl:            "https://webhook.example.com",
					},
				},
				"Email": []models.NotificationInfo{
					{
						EmailNotificationsSettings: true,
						Email:                      "test@gmail.com",
					},
				},
			},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				return mApp
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			appRepository := testScenario.setupMock()
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			sortedNotificationsToSend := appNotificationsService.assignNotificationToProperSendService(testScenario.notifications)
			assert.Equal(t, testScenario.expectedSortedNotifications, sortedNotificationsToSend)
		})
	}
}

func TestAppNotificationsService_sortNotificationsBySendToInformation(t *testing.T) {
	type args struct {
		name                         string
		sortedNotifications          map[string][]models.NotificationInfo
		expectedDiscordNotifications map[string]string
		expectedSlackNotifications   map[string]string
		setupMock                    func() interfaces.AppRepository
	}
	testsScenarios := []args{
		{
			name: "Properly sorted notifications by send to information",
			sortedNotifications: map[string][]models.NotificationInfo{
				"Discord": []models.NotificationInfo{
					{
						DiscordNotificationsSettings: true,
						DiscordWebhookUrl:            "https://webhook.example.discord.com",
						ID:                           "Discord",
						Name:                         "Discord",
						Status:                       "stopped",
					},
				},
				"Slack": []models.NotificationInfo{
					{
						SlackNotificationsSettings: true,
						SlackWebhookUrl:            "https://webhook.example.slack.com",
						ID:                         "Slack",
						Name:                       "Slack",
						Status:                     "stopped",
					},
				},
			},
			expectedDiscordNotifications: map[string]string{
				"https://webhook.example.discord.com": "Discord - Discord - stopped\n",
			},
			expectedSlackNotifications: map[string]string{
				"https://webhook.example.slack.com": "Slack - Slack - stopped\n",
			},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				return mApp
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			appRepository := testScenario.setupMock()
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			discordNotifications, slackNotifications := appNotificationsService.sortNotificationsBySendToInformation(testScenario.sortedNotifications)
			assert.Equal(t, testScenario.expectedDiscordNotifications, discordNotifications)
			assert.Equal(t, testScenario.expectedSlackNotifications, slackNotifications)
		})
	}
}

func TestAppNotificationsService_sendWebhook(t *testing.T) {
	type args struct {
		name                string
		notificationsToSend map[string]string
		expectedError       error
		setupMock           func() interfaces.AppRepository
	}
	testsScenarios := []args{
		{
			name: "Failed to send discord notification",
			notificationsToSend: map[string]string{
				"https://webhook.example.discord.com": "Discord - Discord - stopped\n",
			},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				return mApp
			},
			expectedError: errors.New("Post \"https://webhook.example.discord.com\": dial tcp: lookup webhook.example.discord.com: no such host"),
		},
		{
			name: "Failed to send slack notification",
			notificationsToSend: map[string]string{
				"https://webhook.example.slack.com": "Slack - Slack - stopped\n",
			},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				return mApp
			},
			expectedError: errors.New("Post \"https://webhook.example.slack.com\": tls: failed to verify certificate: x509: certificate is valid for *.slack.com, slack.com, not webhook.example.slack.com"),
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			appRepository := testScenario.setupMock()
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			err := appNotificationsService.sendWebhook(ctx, testScenario.notificationsToSend)
			assert.Equal(t, testScenario.expectedError.Error(), err.Error())
		})
	}
}

func TestAppNotificationsService_SendNotifications(t *testing.T) {
	type args struct {
		name          string
		expectedError error
		appsStatuses  []DTO.AppStatus
		setupMock     func() interfaces.AppRepository
	}
	testsScenarios := []args{
		{
			name:          "No app statuses",
			expectedError: nil,
			appsStatuses:  []DTO.AppStatus{},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				return mApp
			},
		},
		{
			name:          "failed to get users to send notifications",
			expectedError: errors.New("failed to get users to send notifications"),
			appsStatuses:  []DTO.AppStatus{{AppID: "32"}},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetUsersToSendNotifications", mock.Anything,
					mock.Anything).Return([]models.NotificationInfo{}, errors.New("failed to get users to send notifications"))
				return mApp
			},
		},
		{
			name:          "Proper data",
			expectedError: errors.New("Post \"https://webhook.example.slack.com\": tls: failed to verify certificate: x509: certificate is valid for *.slack.com, slack.com, not webhook.example.slack.com"),
			appsStatuses:  []DTO.AppStatus{{AppID: "32", Status: "running"}},
			setupMock: func() interfaces.AppRepository {
				mApp := new(mocks.MockAppRepository)
				mApp.On("GetUsersToSendNotifications", mock.Anything,
					mock.Anything).Return([]models.NotificationInfo{
					{
						Status:                       "running",
						SlackNotificationsSettings:   true,
						DiscordNotificationsSettings: true,
						DiscordWebhookUrl:            "",
						SlackWebhookUrl:              "https://webhook.example.slack.com",
					},
				}, nil)
				return mApp
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			loggerService := tests.CreateLogger()
			appRepository := testScenario.setupMock()
			appNotificationsService := NewAppNotificationsService(appRepository, loggerService)
			err := appNotificationsService.SendNotifications(ctx,
				testScenario.appsStatuses)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}
