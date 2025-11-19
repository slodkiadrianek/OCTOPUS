package controllers

import (
	// "context"

	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type userService interface {
	GetUser(ctx context.Context, userId int) (models.User, error)
	InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error
	UpdateUser(ctx context.Context, user DTO.CreateUser, userId int) error
	UpdateUserNotifications(ctx context.Context, userId int, userNotifications DTO.UpdateUserNotificationsSettings) error
	DeleteUser(ctx context.Context, userId int, password string) error
	ChangeUserPassword(ctx context.Context, userId int, currentPassword string, newPassword string) error
}
type UserController struct {
	userService   userService
	loggerService utils.LoggerService
}

func NewUserController(userService userService, loggerService utils.LoggerService) *UserController {
	return &UserController{
		userService:   userService,
		loggerService: loggerService,
	}
}

func (u *UserController) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userIdString, err := utils.ReadParam(r, "userId")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	userIdFromJwt, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	err = utils.ValidateUsersIds(userId, userIdFromJwt)
	if err != nil {
		u.loggerService.Error("You are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIdToken": userIdFromJwt,
		})
		utils.SetError(w, r, err)
		return
	}
	user, err := u.userService.GetUser(r.Context(), userId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 200, user)
}

func (u *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[DTO.CreateUser](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		utils.SetError(w, r, err)
		return
	}

	userDto := DTO.NewCreateUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.userService.InsertUserToDb(r.Context(), *userDto, userBody.Password)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 201, map[string]string{})
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[DTO.UpdateUser](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		utils.SetError(w, r, err)
		return
	}

	userIdString, err := utils.ReadParam(r, "userId")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	userIdFromJwt, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	err = utils.ValidateUsersIds(userId, userIdFromJwt)
	if err != nil {
		u.loggerService.Error("You are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIdToken": userIdFromJwt,
		})
		utils.SetError(w, r, err)
		return
	}
	userDto := DTO.NewCreateUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.userService.UpdateUser(r.Context(), *userDto, userId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]string{})
}
func (u *UserController) UpdateUserNotifications(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[DTO.UpdateUserNotificationsSettings](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		utils.SetError(w, r, err)
		return
	}

	userIdString, err := utils.ReadParam(r, "userId")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	userIdFromJwt, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	err = utils.ValidateUsersIds(userId, userIdFromJwt)
	if err != nil {
		u.loggerService.Error("You are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIdToken": userIdFromJwt,
		})
		utils.SetError(w, r, err)
		return
	}
	err = u.userService.UpdateUserNotifications(r.Context(), userId, *userBody)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]string{})
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[DTO.DeleteUser](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		utils.SetError(w, r, err)
		return
	}

	userIdString, err := utils.ReadParam(r, "userId")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	userIdFromJwt, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	err = utils.ValidateUsersIds(userId, userIdFromJwt)
	if err != nil {
		u.loggerService.Error("You are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIdToken": userIdFromJwt,
		})
		utils.SetError(w, r, err)
		return
	}
	err = u.userService.DeleteUser(r.Context(), userId, userBody.Password)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]string{})
}

func (u *UserController) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[DTO.ChangeUserPassword](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		utils.SetError(w, r, err)
		return
	}

	userIdString, err := utils.ReadParam(r, "userId")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	userIdFromJwt, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	err = utils.ValidateUsersIds(userId, userIdFromJwt)
	if err != nil {
		u.loggerService.Error("You are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIdToken": userIdFromJwt,
		})
		utils.SetError(w, r, err)
		return
	}
	if userBody.NewPassword != userBody.ConfirmPassword {
		err = fmt.Errorf("passwords do not match")
		utils.SetError(w, r, err)
		return
	}

	err = u.userService.ChangeUserPassword(r.Context(), userId, userBody.CurrentPassword, userBody.NewPassword)
	if err != nil {
		utils.SetError(w, r, err)
	}

	utils.SendResponse(w, 204, map[string]string{})

}
