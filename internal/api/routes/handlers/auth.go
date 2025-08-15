package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
)

type AuthHandlers struct {
	UserController *controllers.UserController
}

func NewAuthHandler(userController *controllers.UserController) *AuthHandlers {
	return &AuthHandlers{
		UserController: userController,
	}
}

func (a *AuthHandlers) SetupAuthHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/auth")
	groupRouter.POST("/register", a.UserController)
	groupRouter.POST("/login", )
	groupRouter.GET("/check")
}
