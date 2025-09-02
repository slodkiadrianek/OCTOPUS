package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	LoggerService  *logger.Logger
	UserRepository *repository.UserRepository
}

func NewUserService(loggerService *logger.Logger, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		LoggerService:  loggerService,
		UserRepository: userRepository,
	}
}

func (u *UserService) InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error {
	doesUserExists, err := u.UserRepository.FindUserByEmail(ctx, user.Email)
	if err != nil && err.Error() != "User not found" {
		return err
	}
	if doesUserExists > 0 {
		u.LoggerService.Info("User with this email already exists", user.Email)
		return models.NewError(400, "Verification", "User with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.LoggerService.Info("failed to generate password", err)
		return err
	}
	err = u.UserRepository.InsertUserToDb(ctx, user, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, user DTO.CreateUser, userId int) error {
	err := u.UserRepository.UpdateUser(ctx, user, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, userId int, password string) error {
	err := u.UserRepository.DeleteUser(ctx, password, userId)
	if err != nil {
		return err
	}
	return nil
}
