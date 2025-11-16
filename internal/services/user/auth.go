package user

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	LoggerService  *utils.Logger
	UserRepository interfaces.UserRepository
	JWT            *middleware.JWT
}

func NewAuthService(loggerService *utils.Logger, userRepository interfaces.UserRepository,
	jwt *middleware.JWT) *AuthService {
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
	if user.ID == 0 {
		a.LoggerService.Info("User with this email does not exist", loginData.Email)
		return "", models.NewError(400, "Verification", "User with this email does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		a.LoggerService.Info("Wrong password provided", loginData)
		return "", models.NewError(401, "Authorization", "Wrong user provided")
	}
	loggedUser := DTO.NewLoggedUser(user.ID, user.Email, user.Name, user.Surname)
	authorizationToken, err := a.JWT.GenerateToken(*loggedUser)
	return authorizationToken, nil
}
