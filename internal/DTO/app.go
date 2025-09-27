package DTO

import "time"

type App struct {
	Id          string
	Name        string
	Description string
	IsDocker    bool
	OwnerID     int
	IpAddress   string
	Port        string
}

func NewApp(id, name, description string, isDocker bool, ownerId int, ipAddress, port string) *App {
	return &App{
		Id:          id,
		Name:        name,
		Description: description,
		IsDocker:    isDocker,
		OwnerID:     ownerId,
		IpAddress:   ipAddress,
		Port:        port,
	}
}

type CreateApp struct {
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app"`
	IpAddress      string `json:"ipAddress" example:"192.168.0.100"`
	Port           string `json:"port" example:"3030"`
	DiscordWebhook string `json:"discordWebhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhook   string `json:"slackWebhook" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}
type UpdateApp struct {
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app"`
	IpAddress      string `json:"ipAddress" example:"192.168.0.100"`
	Port           string `json:"port" example:"3030"`
	DiscordWebhook string `json:"discordWebhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhook   string `json:"slackWebhook" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

type AppId struct {
	AppId string `json:"appId" example:"nd3289dh23934382"`
}
type AppStatus struct {
	AppId     string        `json:"app_id"`
	Status    string        `json:"status"`
	ChangedAt time.Time     `json:"changed_at"`
	Duration  time.Duration `json:"duration"`
}

func NewAppStatus(appId, status string, changedAt time.Time, duration time.Duration) *AppStatus {
	return &AppStatus{
		AppId:     appId,
		Status:    status,
		ChangedAt: changedAt,
		Duration:  duration,
	}
}
