package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
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
		a.loggerService.Error(failedToReadBodyFromRequest, err)
		utils.SetError(w, r, err)
		return
	}

	tokenString, err := a.authService.LoginUser(r.Context(), *userBody)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 200, map[string]string{"token": tokenString})
}

func (a AuthController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]int{
		"userId": userId,
	})
}

func (a AuthController) LogoutUser(w http.ResponseWriter, _ *http.Request) {
	utils.SendResponse(w, 204, map[string]string{})
}
