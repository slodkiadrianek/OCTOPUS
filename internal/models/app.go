package models

type App struct {
	Id             string `json:"id" example:"1"`
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app."`
	IsDocker       bool   `json:"is_docker" example:"false"`
	OwnerID        int    `json:"owner_id" example:"1"`
	IpAddress      string `json:"ip_address" example:"192.168.1.1"`
	Port           string `json:"port" example:"8080"`
	SlackWebhook   string `json:"slack_webhook" example:"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"`
	DiscordWebhook string `json:"discord_webhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

type AppToCheck struct {
	Id        string `json:"id" example:"1"`
	Name      string `json:"name" example:"My App"`
	OwnerID   int    `json:"owner_id" example:"1"`
	IsDocker  bool   `json:"is_docker" example:"false"`
	IpAddress string `json:"ip_address" example:"192.168.1.1"`
	Port      string `json:"port" example:"8080"`
	Status    string `json:"status" example:"running"`
}

type SendNotificationTo struct {
	Id                   string `json:"id" example:"1"`
	Name                 string `json:"name" example:"My App"`
	Email                string `json:"email" sql:"email" example:"joedoe@email.com"`
	Status               string `json:"status" example:"running"`
	SlackWebhook         string `json:"slack_webhook" example:"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"`
	DiscordWebhook       string `json:"discord_webhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	DiscordNotifications bool   `json:"discordNotifications" sql:"discord_notifications" example:"false"`
	EmailNotifications   bool   `json:"emailNotifications" sql:"email_notifications" example:"true"`
	SlackNotifications   bool   `json:"slackNotifications" sql:"slack_notifications" example:"false"`
}

func NewAppToCheck(id, name string, ownerID int, isDocker bool, ipAddress, port string, status string) *AppToCheck {
	return &AppToCheck{
		Id:        id,
		Name:      name,
		OwnerID:   ownerID,
		IsDocker:  isDocker,
		IpAddress: ipAddress,
		Port:      port,
		Status:    status,
	}
}

// func NewApp(name, description, dbLink, apiLink, slackWebhook, discordWebhook string, id, ownerID int) *App {
// 	return &App{
// 		Id:             id,
// 		Name:           name,
// 		Description:    description,
// 		DbLink:         dbLink,
// 		ApiLink:        apiLink,
// 		OwnerID:        ownerID,
// 		SlackWebhook:   slackWebhook,
// 		DiscordWebhook: discordWebhook,
// 	}
// }
