package models

type App struct {
	ID                string `json:"id" example:"1"`
	Name              string `json:"name" example:"My App"`
	Description       string `json:"description" example:"This is my app."`
	IsDocker          bool   `json:"is_docker" example:"false"`
	OwnerID           int    `json:"owner_id" example:"1"`
	IPAddress         string `json:"ip_address" example:"192.168.1.1"`
	Port              string `json:"port" example:"8080"`
	SlackWebhookURL   string `json:"slack_webhook_url" example:"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"`
	DiscordWebhookURL string `json:"discord_webhook_url" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

type AppToCheck struct {
	ID        string `json:"id" example:"1"`
	Name      string `json:"name" example:"My App"`
	OwnerID   int    `json:"owner_id" example:"1"`
	IsDocker  bool   `json:"is_docker" example:"false"`
	IPAddress string `json:"ip_address" example:"192.168.1.1"`
	Port      string `json:"port" example:"8080"`
	Status    string `json:"status" example:"running"`
}

type NotificationInfo struct {
	ID                           string `json:"id" example:"1"`
	Name                         string `json:"name" example:"My App"`
	Email                        string `json:"email" sql:"email" example:"joedoe@email.com"`
	Status                       string `json:"status" example:"running"`
	SlackWebhookURL              string `json:"slack_webhook_url" example:"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"`
	DiscordWebhookURL            string `json:"discord_webhook_url" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	DiscordNotificationsSettings bool   `json:"discordNotificationsSettings " sql:"discord_notifications" example:"false"`
	EmailNotificationsSettings   bool   `json:"emailNotificationsSettings " sql:"email_notifications" example:"true"`
	SlackNotificationsSettings   bool   `json:"slackNotificationsSettings " sql:"slack_notifications" example:"false"`
}
