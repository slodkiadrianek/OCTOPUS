package models

import "time"

type User struct {
	ID                           int       `json:"id" sql:"id" example:"1"`
	Email                        string    `json:"email" sql:"email" example:"joedoe@email.com"`
	Name                         string    `json:"name" sql:"name" example:"Joe"`
	Surname                      string    `json:"surname" sql:"surname" example:"Doe"`
	Password                     string    `json:"password" example:"fsdf2332@!32"`
	DiscordNotificationsSettings bool      `json:"discord_notifications" sql:"discord_notifications" example:"false"`
	EmailNotificationsSettings   bool      `json:"email_notifications_settings" sql:"email_notifications" example:"true"`
	SlackNotificationsSettings   bool      `json:"slack_notifications_settings" sql:"slack_notifications" example:"false"`
	CreatedAt                    time.Time `json:"created_at" sql:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt                    time.Time `json:"updated_at" sql:"updated_at" example:"2023-01-01T00:00:00Z"`
}
