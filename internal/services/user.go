package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
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

func (u *UserService) InsertUserToDb(ctx context.Context, user DTO.User, password string) error {
	err := u.UserRepository.InsertUserToDb(ctx, user, password)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, user DTO.User, userId int) error {
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
