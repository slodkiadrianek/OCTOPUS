package schema

import z "github.com/Oudwins/zog"

var CreateAppSchema = z.Struct(z.Shape{
	"name":              z.String().Required(),
	"description":       z.String().Optional(),
	"ipAddress":         z.String().Required(),
	"port":              z.String().Required(),
	"discordWebhookUrl": z.String().Optional(),
	"slackWebhookUrl":   z.String().Optional(),
})

var AppIdSchema = z.Struct(z.Shape{
	"appId": z.String().Required(),
})
var UpdateAppSchema = z.Struct(z.Shape{
	"name":              z.String().Required(),
	"description":       z.String().Optional(),
	"ipAddress":         z.String().Required(),
	"port":              z.String().Required(),
	"discordWebhookUrl": z.String().Optional(),
	"slackWebhookUrl":   z.String().Optional(),
})
