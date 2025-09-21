package schema

import z "github.com/Oudwins/zog"

type CreateApp struct {
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app"`
	IpAddress      string `json:"ipAddress" example:"192.168.0.100"`
	Port           string `json:"port" example:"3030"`
	DiscordWebhook string `json:"discordWebhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhook   string `json:"slackWebhook" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

var CreateAppSchema = z.Struct(z.Shape{
	"name":           z.String().Required(),
	"description":    z.String().Optional(),
	"ipAddress":      z.String().Required(),
	"port":           z.String().Required(),
	"discordWebhook": z.String().Optional(),
	"slackWebhook":   z.String().Optional(),
})

type AppId struct {
	AppId string `json:"appId" example:"nd3289dh23934382"`
}

var AppIdSchema = z.Struct(z.Shape{
	"appId": z.String().Required(),
})

type UpdateApp struct {
	Name           string `json:"name" example:"My App"`
	Description    string `json:"description" example:"This is my app"`
	IpAddress      string `json:"ipAddress" example:"192.168.0.100"`
	Port           string `json:"port" example:"3030"`
	DiscordWebhook string `json:"discordWebhook" example:"https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz"`
	SlackWebhook   string `json:"slackWebhook" example:"https://hooks.slack.com/services/1234567890/abcdefghijklmnopqrstuvwxyz"`
}

var UpdateAppSchema = z.Struct(z.Shape{
	"name":           z.String().Required(),
	"description":    z.String().Optional(),
	"ipAddress":      z.String().Required(),
	"port":           z.String().Required(),
	"discordWebhook": z.String().Optional(),
	"slackWebhook":   z.String().Optional(),
})
