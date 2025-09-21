package models

import "time"

type User struct {
	Id                   int       `json:"id" sql:"id" example:"1"`
	Email                string    `json:"email" sql:"email" example:"joedoe@email.com"`
	Name                 string    `json:"name" sql:"name" example:"Joe"`
	Surname              string    `json:"surname" sql:"surname" example:"Doe"`
	Password             string    `json:"password" example:"fsdf2332@!32"`
	DiscordNotifications bool      `json:"discordNotifications" sql:"discord_notifications" example:"false"`
	EmailNotifications   bool      `json:"emailNotifications" sql:"email_notifications" example:"true"`
	SlackNotifications   bool      `json:"slackNotifications" sql:"slack_notifications" example:"false"`
	CreatedAt            time.Time `json:"createdAt" sql:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt            time.Time `json:"updatedAt" sql:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// var query string = "CREATE TABLE IF NOT EXISTS users (" +
//	"id INT PRIMARY KEY AUTOINCREMENT," +
//	"email VARCHAR(128) UNIQUE," +
//	"name VARCHAR(64)," +
//	"surname VARCHAR(64)," +
//	"role VARCHAR(64)," +
//	")"
