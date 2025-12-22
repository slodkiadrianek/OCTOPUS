package interfaces

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
	InsertUserToDB(ctx context.Context, user DTO.CreateUser, password string) error
	UpdateUser(ctx context.Context, user DTO.CreateUser, userID int) error
	UpdateUserNotifications(ctx context.Context, userID int, userNotifications DTO.UpdateUserNotificationsSettings) error
	DeleteUser(ctx context.Context, password string, userID int) error
	FindUserByID(ctx context.Context, userID int) (models.User, error)
	ChangeUserPassword(ctx context.Context, userID int, newPassword string) error
}
