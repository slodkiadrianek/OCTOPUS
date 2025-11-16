package interfaces

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
	InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error
	UpdateUser(ctx context.Context, user DTO.CreateUser, userId int) error
	UpdateUserNotifications(ctx context.Context, userId int, userNotifications DTO.UpdateUserNotificationsSettings) error
	DeleteUser(ctx context.Context, password string, userId int) error
	FindUserById(ctx context.Context, userId int) (models.User, error)
	ChangeUserPassword(ctx context.Context, userId int, newPassword string) error
}
