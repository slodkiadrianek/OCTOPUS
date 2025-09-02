package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/middleware"
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
		JWT:         jwt,
	}
}

func (a AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.LoginUser](r)
	if err != nil {
		return
	}
}
