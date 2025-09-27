package schema

import z "github.com/Oudwins/zog"

var CreateAppSchema = z.Struct(z.Shape{
	"name":           z.String().Required(),
	"description":    z.String().Optional(),
	"ipAddress":      z.String().Required(),
	"port":           z.String().Required(),
	"discordWebhook": z.String().Optional(),
	"slackWebhook":   z.String().Optional(),
})

var AppIdSchema = z.Struct(z.Shape{
	"appId": z.String().Required(),
})
var UpdateAppSchema = z.Struct(z.Shape{
	"name":           z.String().Required(),
	"description":    z.String().Optional(),
	"ipAddress":      z.String().Required(),
	"port":           z.String().Required(),
	"discordWebhook": z.String().Optional(),
	"slackWebhook":   z.String().Optional(),
})
