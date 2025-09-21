package schema

type CreateApp struct {
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app"`
	DbLink         string `json:"dbLink" example:"mongodb://localhost:27017"`
	ApiLink        string `json:"apiLink" example:"https://api.example.com"`
	DiscordWebhook string `json:"discordWebhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhook   string `json:"slackWebhook" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

type UpdateApp struct {
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app"`
	DbLink         string `json:"dbLink" example:"mongodb://localhost:27017"`
	ApiLink        string `json:"apiLink" example:"https://api.example.com"`
	DiscordWebhook string `json:"discordWebhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhook   string `json:"slackWebhook" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
	ID             int    `json:"id" example:"1"`
}
