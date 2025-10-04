package services

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	LoggerService  *utils.Logger
	UserRepository userRepository
	JWT            *middleware.JWT
}

func NewAuthService(loggerService *utils.Logger, userRepository userRepository, jwt *middleware.JWT) *AuthService {
	return &AuthService{
		LoggerService:  loggerService,
		UserRepository: userRepository,
		JWT:            jwt,
	}
}

func (a AuthService) LoginUser(ctx context.Context, loginData DTO.LoginUser) (string, error) {
	user, err := a.UserRepository.FindUserByEmail(ctx, loginData.Email)
	if err != nil && err.Error() != "User not found" {
		return "", err
	}
	if user.Id == 0 {
		a.LoggerService.Info("User with this email does not exist", loginData.Email)
		return "", models.NewError(400, "Verification", "User with this email does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		a.LoggerService.Info("Wrong password provided", loginData)
		return "", models.NewError(401, "Authorization", "Wron guser provided")
	}
	loggedUser := DTO.NewLoggedUser(user.Id, user.Email, user.Name, user.Surname)
	tokenString, err := a.JWT.GenerateToken(*loggedUser)
	return tokenString, nil
}
