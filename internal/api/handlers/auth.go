package handlers

import "github.com/slodkiadrianek/octopus/internal/api/routes"

type AuthHandlers struct {
	AuthControllers
}

func (a *AuthHandlers) SetupAuthHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/auth")
	groupRouter.POST("/register")
}
