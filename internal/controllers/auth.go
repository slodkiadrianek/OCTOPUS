package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService, jwt *middleware.JWT) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (a AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.LoginUser](r)
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
	utils.SendResponse(w, 204, map[string]string{})
}

func (a AuthController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, 204, map[string]string{})
}
