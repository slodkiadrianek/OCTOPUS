package handlers

import (
	"github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type AuthHandlers struct {
	userController interfaces.UserController
	authController interfaces.AuthController
	jwt            *middleware.JWT
	rateLimiter    *middleware.RateLimiter
}

func NewAuthHandler(userController interfaces.UserController, authController interfaces.AuthController, jwt *middleware.JWT,
	rateLimiter *middleware.RateLimiter) *AuthHandlers {
	return &AuthHandlers{
		userController: userController,
		authController: authController,
		jwt:            jwt,
		rateLimiter:    rateLimiter,
	}
}

func (a *AuthHandlers) SetupAuthHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/auth")

	groupRouter.POST("/register", middleware.RateLimiterMiddleware(*a.rateLimiter),
		middleware.ValidateMiddleware[DTO.CreateUser, *zog.StructSchema](
			"body",
			schema.CreateUserSchema),
		a.userController.InsertUser)
	groupRouter.POST("/login", middleware.ValidateMiddleware[DTO.LoginUser]("body", schema.LoginUserSchema),
		a.authController.LoginUser)

	groupRouter.GET("/check", a.jwt.VerifyToken, a.authController.VerifyUser)
	groupRouter.DELETE("/logout", a.jwt.VerifyToken, a.jwt.BlacklistUser, a.authController.LogoutUser)
}
