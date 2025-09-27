package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	LoggerService  *utils.Logger
	UserRepository *repository.UserRepository
}

func NewUserService(loggerService *utils.Logger, userRepository *repository.UserRepository) *UserService {
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
	if doesUserExists.Id > 0 {
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

func (u *UserService) UpdateUserNotifications(ctx context.Context, userId int,
	userNotifications DTO.UpdateUserNotifications) error {
	err := u.UserRepository.UpdateUserNotifications(ctx, userId, userNotifications)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, userId int, password string) error {
	user, err := u.UserRepository.FindUserById(ctx, userId)
	if err != nil {
		return err
	}
	if user.Id == 0 {
		u.LoggerService.Info("User with this id does not exist", userId)
		return models.NewError(400, "Verification", "User with this id does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		u.LoggerService.Info("Wrong password provided", userId)
		return models.NewError(401, "Authorization", "Wrong password provided")
	}
	err = u.UserRepository.DeleteUser(ctx, password, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) ChangeUserPassword(ctx context.Context, userId int, currentPassword string, newPassword string) error {
	user, err := u.UserRepository.FindUserById(ctx, userId)
	if err != nil {
		return err
	}
	if user.Id == 0 {
		u.LoggerService.Info("User with this id does not exist", userId)
		return models.NewError(400, "Verification", "User with this id does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		u.LoggerService.Info("Wrong current password provided", userId)
		return models.NewError(401, "Authorization", "Wrong current password provided")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		u.LoggerService.Info("failed to generate password", err)
		return err
	}
	err = u.UserRepository.ChangeUserPassword(ctx, userId, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}
