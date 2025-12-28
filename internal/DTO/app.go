package DTO

import "time"

type App struct {
	ID          string
	Name        string
	Description string
	IsDocker    bool
	OwnerID     int
	IPAddress   string
	Port        string
}

func NewApp(id, name, description string, isDocker bool, ownerID int, IPAddress, port string) *App {
	return &App{
		ID:          id,
		Name:        name,
		Description: description,
		IsDocker:    isDocker,
		OwnerID:     ownerID,
		IPAddress:   IPAddress,
		Port:        port,
	}
}

type CreateApp struct {
	Name              string `json:"name" example:"My App"`
	Description       string `json:"description" example:"This is my app"`
	IPAddress         string `json:"ipAddress" example:"192.168.0.100"`
	Port              string `json:"port" example:"3030"`
	DiscordWebhookURL string `json:"discordWebhookUrl" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhookURL   string `json:"slackWebhookUrl" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}
type UpdateApp struct {
	Name              string `json:"name" example:"My App"`
	Description       string `json:"description" example:"This is my app"`
	IPAddress         string `json:"ipAddress" example:"192.168.0.100"`
	Port              string `json:"port" example:"3030"`
	DiscordWebhookURL string `json:"discordWebhookUrl" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhookURL   string `json:"slackWebhookUrl" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

type AppID struct {
	AppID string `json:"appID" example:"nd3289dh23934382"`
}
type AppStatus struct {
	AppID     string        `json:"app_id"`
	Status    string        `json:"status"`
	ChangedAt time.Time     `json:"changed_at"`
	Duration  time.Duration `json:"duration"`
}

func NewAppStatus(appID, status string, changedAt time.Time, duration time.Duration) *AppStatus {
	return &AppStatus{
		AppID:     appID,
		Status:    status,
		ChangedAt: changedAt,
		Duration:  duration,
	}
}
