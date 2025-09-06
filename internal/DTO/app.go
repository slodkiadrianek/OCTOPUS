package DTO

type App struct {
	Name           string `json:"name"`
<<<<<<< HEAD
	Description    string `json:"description"`
=======
>>>>>>> a4f4bb342f74a1e297be363d81262025c784bffa
	DbLink         string `json:"dbLink"`
	ApiLink        string `json:"apiLink"`
	OwnerID        int    `json:"ownerId"`
	DiscordWebhook string `json:"discordWebhook"`
	SlackWebhook   string `json:"slackWebhook"`
}

<<<<<<< HEAD
type UpdateApp struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DbLink         string `json:"dbLink"`
	ApiLink        string `json:"apiLink"`
	DiscordWebhook string `json:"discordWebhook"`
	SlackWebhook   string `json:"slackWebhook"`
}

func NewUpdateApp(id int, name string, description string, dbLink string, apiLink string, discordWebhook string, slackWebhook string) *UpdateApp {
	return &UpdateApp{
		Id:             id,
		Name:           name,
		Description:    description,
		DbLink:         dbLink,
		ApiLink:        apiLink,
		DiscordWebhook: discordWebhook,
		SlackWebhook:   slackWebhook,
	}
}

func NewApp(name string, description string, dbLink string, apiLink string, ownerId int, discordWebhook string, slackWebhook string) *App {
	return &App{
		Name:           name,
		Description:    description,
=======
func NewApp(name string, dbLink string, apiLink string, ownerId int, discordWebhook string, slackWebhook string) *App {
	return &App{
		Name:           name,
>>>>>>> a4f4bb342f74a1e297be363d81262025c784bffa
		DbLink:         dbLink,
		ApiLink:        apiLink,
		OwnerID:        ownerId,
		DiscordWebhook: discordWebhook,
		SlackWebhook:   slackWebhook,
	}
}
