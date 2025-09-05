package models

type App struct {
	Id             int    `json:"id" example:"1"`
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app."`
	DbLink         string `json:"db_link" example:"mongodb://localhost:27017"`
	ApiLink        string `json:"api_link" example:"http://localhost:8080"`
	OwnerID        int    `json:"owner_id" example:"1"`
	SlackWebhook   string `json:"slack_webhook" example:"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"`
	DiscordWebhook string `json:"discord_webhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

func NewApp(name, description, dbLink, apiLink, slackWebhook, discordWebhook string, id, ownerID int) *App {
	return &App{
		Id:             id,
		Name:           name,
		Description:    description,
		DbLink:         dbLink,
		ApiLink:        apiLink,
		OwnerID:        ownerID,
		SlackWebhook:   slackWebhook,
		DiscordWebhook: discordWebhook,
	}
}
