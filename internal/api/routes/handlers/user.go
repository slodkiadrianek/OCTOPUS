package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type UserHandlers struct {
	UserController *controllers.UserController
	JWT            *middleware.JWT
}

func NewUserHandler(userController *controllers.UserController, jwt *middleware.JWT) *UserHandlers {
	return &UserHandlers{
		UserController: userController,
		JWT:            jwt,
	}
}

func (u *UserHandlers) SetupUserHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/user")
	groupRouter.PUT("/:userId", middleware.ValidateMiddleware[schema.UpdateUser]("body", schema.UpdateUserSchema), u.UserController.UpdateUser)
	groupRouter.PUT("/:userId/notifications", middleware.ValidateMiddleware[schema.UpdateUserNotifications]("body", schema.UpdateUserNotificationsSchema), u.UserController.UpdateUserNotifications)
	groupRouter.DELETE("/:userId", middleware.ValidateMiddleware[schema.UserId]("params", schema.UserIdSchema), u.JWT.BlacklistUser, u.UserController.DeleteUser)
}
