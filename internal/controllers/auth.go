package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
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
	userBody, err := request.ReadBody[DTO.LoginUser](r)
	if err != nil {
		a.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	tokenString, err := a.authService.LoginUser(r.Context(), *userBody)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, map[string]string{"token": tokenString})
}

func (a AuthController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	userId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]int{
		"userId": userId,
	})
}

func (a AuthController) LogoutUser(w http.ResponseWriter, _ *http.Request) {
	response.Send(w, 204, map[string]string{})
}
