package schema

import (
	"strings"

	z "github.com/Oudwins/zog"
)

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

var LoginUserSchema = z.Struct(z.Shape{
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
	"password": z.String().Required(),
})

var UpdateUserSchema = z.Struct(z.Shape{
	"name":    z.String().Required(),
	"surname": z.String().Required(),
	"email": z.String().Email().Required().Transform(func(val *string, ctx z.Ctx) error {
		*val = strings.ToLower(*val)
		*val = strings.TrimSpace(*val)
		return nil
	}),
})

var UserIdSchema = z.Struct(z.Shape{
	"userId": z.String().Required(),
})

var ChangeUserPasswordSchema = z.Struct(z.Shape{
	"newPassword":     z.String().Required(),
	"confirmPassword": z.String().Required(),
	"currentPassword": z.String().Required(),
})

var DeleteUserSchema = z.Struct(z.Shape{
	"password": z.String().Required(),
})

var UpdateUserNotificationsSchema = z.Struct(z.Shape{
	"discordNotificationsSettings": z.Bool().Optional(),
	"slackNotificationsSettings":   z.Bool().Optional(),
	"emailNotificationsSettings":   z.Bool().Optional(),
})
