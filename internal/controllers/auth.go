package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AuthController struct {
	AuthService *services.AuthService
	Logger      *utils.Logger
}

func NewAuthController(authService *services.AuthService, logger *utils.Logger) *AuthController {
	return &AuthController{
		AuthService: authService,
		Logger:      logger,
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
	tokenString, err := a.AuthService.LoginUser(r.Context(), *userBody)
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
	userId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if userId == 0 {
		return
	}
	utils.SendResponse(w, 204, map[string]int{
		"userId": userId,
	})
}

func (a AuthController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, 204, map[string]string{})
}
