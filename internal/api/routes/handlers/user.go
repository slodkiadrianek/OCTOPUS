package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type UserHandlers struct {
	userController interfaces.UserController
	jwt            *middleware.JWT
}

func NewUserHandler(userController interfaces.UserController, jwt *middleware.JWT) *UserHandlers {
	return &UserHandlers{
		userController: userController,
		jwt:            jwt,
	}
}

func (u *UserHandlers) SetupUserHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/users")

	groupRouter.GET("/:userId", u.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.UserID]("params",
		schema.UserIdSchema), u.userController.GetUserInfo)

	groupRouter.PUT("/:userId", u.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.UpdateUser]("body",
		schema.UpdateUserSchema), u.userController.UpdateUser)
	groupRouter.PUT("/:userId/notifications", u.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.UserID]("params",
		schema.UserIdSchema), middleware.ValidateMiddleware[DTO.UpdateUserNotificationsSettings]("body",
		schema.UpdateUserNotificationsSchema), u.userController.UpdateUserNotifications)

	groupRouter.PATCH("/:userId", u.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.UserID]("params",
		schema.UserIdSchema), middleware.ValidateMiddleware[DTO.ChangeUserPassword]("body",
		schema.ChangeUserPasswordSchema), u.userController.ChangeUserPassword)
	groupRouter.DELETE("/:userId", u.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.UserID]("params",
		schema.UserIdSchema), middleware.ValidateMiddleware[DTO.DeleteUser]("body", schema.DeleteUserSchema),
		u.userController.DeleteUser)
}
