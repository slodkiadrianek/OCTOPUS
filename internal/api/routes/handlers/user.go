package handlers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type userController interface {
	GetUserInfo(w http.ResponseWriter, r *http.Request)
	InsertUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	UpdateUserNotifications(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ChangeUserPassword(w http.ResponseWriter, r *http.Request)
}
type UserHandlers struct {
	UserController userController
	JWT            *middleware.JWT
}

func NewUserHandler(userController userController, jwt *middleware.JWT) *UserHandlers {
	return &UserHandlers{
		UserController: userController,
		JWT:            jwt,
	}
}

func (u *UserHandlers) SetupUserHandlers(router routes.Router) {
	groupRouter := router.Group("/api/v1/users")
	groupRouter.GET("/:userId", u.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.UserId]("params",
		schema.UserIdSchema), u.UserController.GetUserInfo)
	groupRouter.PUT("/:userId", u.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.UpdateUser]("body",
		schema.UpdateUserSchema), u.UserController.UpdateUser)
	groupRouter.PUT("/:userId/notifications", u.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.UserId]("params",
		schema.UserIdSchema), middleware.ValidateMiddleware[DTO.UpdateUserNotifications]("body",
		schema.UpdateUserNotificationsSchema), u.UserController.UpdateUserNotifications)
	groupRouter.PATCH("/:userId", u.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.UserId]("params",
		schema.UserIdSchema), middleware.ValidateMiddleware[DTO.ChangeUserPassword]("body",
		schema.ChangeUserPasswordSchema), u.UserController.ChangeUserPassword)
	groupRouter.DELETE("/:userId", u.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.UserId]("params",
		schema.UserIdSchema), middleware.ValidateMiddleware[DTO.DeleteUser]("body", schema.DeleteUserSchema),
		u.UserController.DeleteUser)
}
