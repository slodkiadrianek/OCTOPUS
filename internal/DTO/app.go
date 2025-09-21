package DTO

import "time"

type App struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDocker    bool   `json:"is_docker"`
	OwnerID     int    `json:"owner_id"`
	IpAddress   string `json:"ip_address"`
	Port        string `json:"port"`
}

type UpdateApp struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ApiLink        string `json:"apiLink"`
	DiscordWebhook string `json:"discordWebhook"`
	SlackWebhook   string `json:"slackWebhook"`
}

type AppStatus struct {
	AppId     string        `json:"app_id"`
	Status    string        `json:"status"`
	ChangedAt time.Time     `json:"changed_at"`
	Duration  time.Duration `json:"duration"`
}

func NewUpdateApp(id int, name string, description string, dbLink string, apiLink string, discordWebhook string, slackWebhook string) *UpdateApp {
	return &UpdateApp{
		Id:             id,
		Name:           name,
		Description:    description,
		ApiLink:        apiLink,
		DiscordWebhook: discordWebhook,
		SlackWebhook:   slackWebhook,
	}
}

func NewApp(id, name, description string, IsDocker bool, ownerId int, ipAddress string, port string) *App {
	return &App{
		Id:          id,
		Name:        name,
		Description: description,
		IsDocker:    IsDocker,
		OwnerID:     ownerId,
		IpAddress:   ipAddress,
		Port:        port,
	}
}

func NewAppStatus(appId, status string, changedAt time.Time, duration time.Duration) *AppStatus {
	return &AppStatus{
		AppId:     appId,
		Status:    status,
		ChangedAt: changedAt,
		Duration:  duration,
	}
}
