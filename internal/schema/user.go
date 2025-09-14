package schema

import (
	"strings"

	z "github.com/Oudwins/zog"
)

type CreateUser struct {
	Name     string `json:"name" example:"Joe"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"joedoe@email.com"`
	Password string `json:"password" example:"2r3c23rc3#@r32rs2"`
}

var CreateUserSchema = z.Struct(z.Shape{
	"name":    z.String().Required(),
	"surname": z.String().Required(),
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
	"password": z.String().Min(8).Max(32).ContainsSpecial().ContainsUpper().ContainsDigit().Required(),
})

type LoginUser struct {
	Email    string `json:"email" example:"adikurek@gmail.com"`
	Password string `json:"password" example:"zaqwekflas;h#&"`
}

var LoginUserSchema = z.Struct(z.Shape{
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
	"password": z.String().Required(),
})

type UpdateUser struct {
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"joedoe@email.com"`
}

var UpdateUserSchema = z.Struct(z.Shape{
	"name":    z.String().Required(),
	"surname": z.String().Required(),
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
})

type UserId struct {
	UserId int `json:"userId" example:"2"`
}

var UserIdSchema = z.Struct(z.Shape{
	"userId": z.Int().Required(),
})

type DeleteUser struct {
	Password string `json:"password" example:"zaqwekflas;h#&"`
	UserId   UserId
}

var DeleteUserSchema = z.Struct(z.Shape{
	"password": z.String().Required(),
}).Merge(UserIdSchema)

type UpdateUserNotifications struct {
	DiscordNotifications bool `json:"discordNotifications" example:"true"`
	SlackNotifications   bool `json:"slackNotifications" example:"true"`
	EmailNotifications   bool `json:"emailNotifications" example:"true"`
}

var UpdateUserNotificationsSchema = z.Struct(z.Shape{
	"discordNotifications": z.Bool().Optional(),
	"slackNotifications":   z.Bool().Optional(),
	"emailNotifications":   z.Bool().Optional(),
})
