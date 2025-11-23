package user

import (
	"context"
	"fmt"
	"time"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	loggerService  utils.LoggerService
	cacheService   interfaces.CacheService
	userRepository interfaces.UserRepository
}

func NewUserService(loggerService utils.LoggerService, userRepository interfaces.UserRepository,
	cacheService interfaces.CacheService) *UserService {
	return &UserService{
		loggerService:  loggerService,
		userRepository: userRepository,
		cacheService:   cacheService,
	}
}
func (u *UserService) readUserFromCache(ctx context.Context, cacheKey string) (models.User, error) {
	userJson, err := u.cacheService.GetData(ctx, cacheKey)
	if err != nil {
		return models.User{}, err
	}
	userPtr, err := utils.UnmarshalData[models.User]([]byte(userJson))
	if err != nil {
		return models.User{}, err
	}
	user := *userPtr
	return user, nil
}
func (u *UserService) callFindUserByIdAndSaveToCache(ctx context.Context, userId int, cacheKey string) (models.User,
	error) {
	user, err := u.userRepository.FindUserById(ctx, userId)
	if err != nil {
		return models.User{}, err
	}
	if user.ID == 0 {
		return user, nil
	}
	userJson, err := utils.MarshalData(user)
	if err != nil {
		return models.User{}, err
	}
	err = u.cacheService.SetData(ctx, cacheKey, string(userJson), time.Minute)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (u *UserService) GetUser(ctx context.Context, userId int) (models.User, error) {
	cacheKey := fmt.Sprintf("users-%d", userId)
	doesUserExists, err := u.cacheService.ExistsData(ctx, cacheKey)
	if err != nil {
		return models.User{}, err
	}
	if doesUserExists > 0 {
		user, err := u.readUserFromCache(ctx, cacheKey)
		if err != nil {
			return models.User{}, err
		}
		return user, nil
	}
	user, err := u.userRepository.FindUserById(ctx, userId)
	userJson, err := utils.MarshalData(user)
	if err != nil {
		return models.User{}, err
	}
	err = u.cacheService.SetData(ctx, cacheKey, string(userJson), time.Minute)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserService) InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error {
	doesUserExists, err := u.userRepository.FindUserByEmail(ctx, user.Email)
	if err != nil {
		if err.Error() == "User not found" {
			u.loggerService.Info("User not found", user.Email)
			return models.NewError(400, "NotFound", "User not found")
		}
		return models.NewError(500, "InternalError", err.Error())
	}
	if doesUserExists.ID > 0 {
		u.loggerService.Info("User with this email already exists", user.Email)
		return models.NewError(400, "Verification", "User with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.loggerService.Info("failed to generate password", err)
		return err
	}
	err = u.userRepository.InsertUserToDb(ctx, user, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, user DTO.CreateUser, userId int) error {
	err := u.userRepository.UpdateUser(ctx, user, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUserNotifications(ctx context.Context, userId int,
	userNotifications DTO.UpdateUserNotificationsSettings,
) error {
	err := u.userRepository.UpdateUserNotifications(ctx, userId, userNotifications)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, userId int, password string) error {
	cacheKey := fmt.Sprintf("users-%d", userId)
	doesUserExists, err := u.cacheService.ExistsData(ctx, cacheKey)
	if err != nil {
		return err
	}
	var user models.User

	if doesUserExists > 0 {
		user, err = u.readUserFromCache(ctx, cacheKey)
		if err != nil {
			return err
		}
	} else {
		user, err = u.callFindUserByIdAndSaveToCache(ctx, userId, cacheKey)
		if err != nil {
			return err
		}
	}
	if user.ID == 0 {
		u.loggerService.Info("User with this id does not exist", userId)
		return models.NewError(400, "Verification", "User with this id does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		u.loggerService.Info("Wrong password provided", userId)
		return models.NewError(401, "Authorization", "Wrong password provided")
	}
	err = u.userRepository.DeleteUser(ctx, password, userId)
	if err != nil {
		return err
	}
	err = u.cacheService.DeleteData(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) ChangeUserPassword(ctx context.Context, userId int, currentPassword string, newPassword string) error {
	cacheKey := fmt.Sprintf("users-%d", userId)
	doesUserExists, err := u.cacheService.ExistsData(ctx, cacheKey)
	if err != nil {
		return err
	}
	var user models.User

	if doesUserExists > 0 {
		user, err = u.readUserFromCache(ctx, cacheKey)
		if err != nil {
			return err
		}
	} else {
		user, err = u.callFindUserByIdAndSaveToCache(ctx, userId, cacheKey)
		fmt.Println(err)
		if err != nil {
			return err
		}
	}
	if user.ID == 0 {
		u.loggerService.Info("User with this id does not exist", userId)
		return models.NewError(400, "Verification", "User with this id does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		u.loggerService.Info("Wrong current password provided", userId)
		return models.NewError(401, "Authorization", "Wrong current password provided")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		u.loggerService.Info("failed to generate password", err)
		return err
	}
	err = u.userRepository.ChangeUserPassword(ctx, userId, string(hashedPassword))
	if err != nil {
		return err
	}
	err = u.cacheService.DeleteData(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
