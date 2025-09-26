package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type AuthHandlers struct {
	UserController *controllers.UserController
	AuthController *controllers.AuthController
	JWT            *middleware.JWT
}

func NewAuthHandler(userController *controllers.UserController, authController *controllers.AuthController, jwt *middleware.JWT) *AuthHandlers {
	return &AuthHandlers{
		UserController: userController,
		AuthController: authController,
		JWT:            jwt,
	}
}

func (a *AuthHandlers) SetupAuthHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/auth")
	groupRouter.POST("/register", middleware.ValidateMiddleware[schema.CreateUser]("body", schema.CreateUserSchema), a.UserController.InsertUser)
	groupRouter.POST("/login", middleware.ValidateMiddleware[schema.LoginUser]("body", schema.LoginUserSchema), a.AuthController.LoginUser)
	groupRouter.GET("/check", a.JWT.VerifyToken, a.AuthController.VerifyUser)
	groupRouter.DELETE("/logout", a.JWT.VerifyToken, a.JWT.BlacklistUser, a.AuthController.LogoutUser)
}
