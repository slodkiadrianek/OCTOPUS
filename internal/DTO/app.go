package DTO

type App struct {
	Name           string `json:"name"`
	DbLink         string `json:"dbLink"`
	ApiLink        string `json:"apiLink"`
	OwnerID        int    `json:"ownerId"`
	DiscordWebhook string `json:"discordWebhook"`
	SlackWebhook   string `json:"slackWebhook"`
}

func NewApp(name string, dbLink string, apiLink string, ownerId int, discordWebhook string, slackWebhook string) *App {
	return &App{
		Name:           name,
		DbLink:         dbLink,
		ApiLink:        apiLink,
		OwnerID:        ownerId,
		DiscordWebhook: discordWebhook,
		SlackWebhook:   slackWebhook,
	}
}
