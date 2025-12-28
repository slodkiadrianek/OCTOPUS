package schema

import (
	"strings"

	z "github.com/Oudwins/zog"
)

var CreateUserSchema = z.Struct(z.Shape{
	"name":    z.String().Required().Max(64),
	"surname": z.String().Required().Max(64),
	"email": z.String().Email().Required().Max(64).Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
	"password": z.String().Min(8).Max(32).ContainsSpecial().ContainsUpper().ContainsDigit().Required(),
})

var LoginUserSchema = z.Struct(z.Shape{
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}).Max(64),
	"password": z.String().Required(),
})

var UpdateUserSchema = z.Struct(z.Shape{
	"name":    z.String().Required().Max(64),
	"surname": z.String().Required().Max(64),
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
})

var UserIDSchema = z.Struct(z.Shape{
	"userID": z.String().Required(),
})

var ChangeUserPasswordSchema = z.Struct(z.Shape{
	"newPassword":     z.String().Min(8).Max(32).ContainsSpecial().ContainsUpper().ContainsDigit().Required(),
	"confirmPassword": z.String().Required().Max(32),
	"currentPassword": z.String().Required().Max(32),
})

var DeleteUserSchema = z.Struct(z.Shape{
	"password": z.String().Required().Max(32),
})

var UpdateUserNotificationsSchema = z.Struct(z.Shape{
	"discordNotificationsSettings": z.Bool().Optional(),
	"slackNotificationsSettings":   z.Bool().Optional(),
	"emailNotificationsSettings":   z.Bool().Optional(),
})
