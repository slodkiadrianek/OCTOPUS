package handlers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type authController interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
	VerifyUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
}

type AuthHandlers struct {
	UserController userController
	AuthController authController
	JWT            *middleware.JWT
	RateLimiter    *middleware.RateLimiter
}

func NewAuthHandler(userController userController, authController authController, jwt *middleware.JWT,
	rateLimiter *middleware.RateLimiter) *AuthHandlers {
	return &AuthHandlers{
		UserController: userController,
		AuthController: authController,
		JWT:            jwt,
		RateLimiter:    rateLimiter,
	}
}

func (a *AuthHandlers) SetupAuthHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/auth")
	groupRouter.POST("/register", middleware.RateLimiterMiddleware(*a.RateLimiter),
		middleware.ValidateMiddleware[DTO.CreateUser](
			"body",
			schema.CreateUserSchema),
		a.UserController.InsertUser)
	groupRouter.POST("/login", middleware.ValidateMiddleware[DTO.LoginUser]("body", schema.LoginUserSchema),
		a.AuthController.LoginUser)
	groupRouter.GET("/check", a.JWT.VerifyToken, a.AuthController.VerifyUser)
	groupRouter.DELETE("/logout", a.JWT.VerifyToken, a.JWT.BlacklistUser, a.AuthController.LogoutUser)
}
