package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
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
	groupRouter.POST("/register",middleware.ValidateMiddleware[schema.CreateUser]("body", schema.CreateUserSchema), a.UserController.InsertUser)
	groupRouter.POST("/login")
	groupRouter.GET("/check")
}
