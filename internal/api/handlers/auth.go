package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
)

type AuthHandlers struct {
	UserController *controllers.UserController
}

func (a *AuthHandlers) SetupAuthHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/auth")
	groupRouter.POST("/register", a.UserController)
}
