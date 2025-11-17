package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type authService interface {
	LoginUser(ctx context.Context, loginData DTO.LoginUser) (string, error)
}

type AuthController struct {
	authService   authService
	loggerService utils.LoggerService
}

func NewAuthController(authService authService, loggerService utils.LoggerService) *AuthController {
	return &AuthController{
		authService:   authService,
		loggerService: loggerService,
	}
}

func (a AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[DTO.LoginUser](r)
	if err != nil {
		errBucket, ok := r.Context().Value("ErrorBucket").(*models.ErrorBucket)
		if ok {
			errBucket.Err = err
			return
		}
		utils.SendResponse(w, 500, map[string]string{
			"errorCategory":    "Server",
			"errorDescription": "Internal server error",
		})
	}
	tokenString, err := a.authService.LoginUser(r.Context(), *userBody)
	if err != nil {
		errBucket, ok := r.Context().Value("ErrorBucket").(*models.ErrorBucket)
		if ok {
			errBucket.Err = err
			return
		}
		utils.SendResponse(w, 500, map[string]string{
			"errorCategory":    "Server",
			"errorDescription": "Internal server error",
		})
	}
	utils.SendResponse(w, 200, map[string]string{"token": tokenString})
}

func (a AuthController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	userId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if userId == 0 {
		return
	}
	utils.SendResponse(w, 204, map[string]int{
		"userId": userId,
	})
}

func (a AuthController) LogoutUser(w http.ResponseWriter, _ *http.Request) {
	utils.SendResponse(w, 204, map[string]string{})
}
