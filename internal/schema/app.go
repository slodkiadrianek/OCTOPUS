package schema

import z "github.com/Oudwins/zog"

var CreateAppSchema = z.Struct(z.Shape{
	"name":              z.String().Required().Max(64),
	"description":       z.String().Optional().Max(256),
	"ipAddress":         z.String().Required().Max(32),
	"port":              z.String().Required().Max(16),
	"discordWebhookUrl": z.String().Optional().Max(256),
	"slackWebhookUrl":   z.String().Optional().Max(256),
})

var AppIDSchema = z.Struct(z.Shape{
	"appId": z.String().Required().Max(64),
})

var UpdateAppSchema = z.Struct(z.Shape{
	"name":              z.String().Required().Max(64),
	"description":       z.String().Optional().Max(256),
	"ipAddress":         z.String().Required().Max(32),
	"port":              z.String().Required().Max(16),
	"discordWebhookUrl": z.String().Optional().Max(256),
	"slackWebhookUrl":   z.String().Optional().Max(256),
})
