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
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
	"github.com/slodkiadrianek/octopus/internal/utils/validation"
)

type userService interface {
	GetUser(ctx context.Context, userID int) (models.User, error)
	InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error
	UpdateUser(ctx context.Context, user DTO.CreateUser, userID int) error
	UpdateUserNotifications(ctx context.Context, userID int, userNotifications DTO.UpdateUserNotificationsSettings) error
	DeleteUser(ctx context.Context, userID int, password string) error
	ChangeUserPassword(ctx context.Context, userID int, currentPassword string, newPassword string) error
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
	userIDString, err := request.ReadParam(r, "userID")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		u.loggerService.Error("failed to covnert string to int", err)
		response.SetError(w, r, err)
		return
	}
	userIDFromJwt, err := request.ReadUserIDFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = validation.ValidateUsersIDs(userID, userIDFromJwt)
	if err != nil {
		u.loggerService.Error("you are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIDToken": userIDFromJwt,
		})
		response.SetError(w, r, err)
		return
	}

	user, err := u.userService.GetUser(r.Context(), userID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, user)
}

func (u *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := request.ReadBody[DTO.CreateUser](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	userDto := DTO.NewCreateUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.userService.InsertUserToDb(r.Context(), *userDto, userBody.Password)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 201, map[string]string{})
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := request.ReadBody[DTO.UpdateUser](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	userIDString, err := request.ReadParam(r, "userID")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	userIDFromJwt, err := request.ReadUserIDFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = validation.ValidateUsersIDs(userID, userIDFromJwt)
	if err != nil {
		u.loggerService.Error("you are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIDToken": userIDFromJwt,
		})
		response.SetError(w, r, err)
		return
	}

	userDto := DTO.NewCreateUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.userService.UpdateUser(r.Context(), *userDto, userID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]string{})
}

func (u *UserController) UpdateUserNotifications(w http.ResponseWriter, r *http.Request) {
	userBody, err := request.ReadBody[DTO.UpdateUserNotificationsSettings](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	userIDString, err := request.ReadParam(r, "userID")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	userIDFromJwt, err := request.ReadUserIDFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = validation.ValidateUsersIDs(userID, userIDFromJwt)
	if err != nil {
		u.loggerService.Error("you are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIDToken": userIDFromJwt,
		})
		response.SetError(w, r, err)
		return
	}

	err = u.userService.UpdateUserNotifications(r.Context(), userID, *userBody)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]string{})
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := request.ReadBody[DTO.DeleteUser](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	userIDString, err := request.ReadParam(r, "userID")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		u.loggerService.Error("failed to covnert string to int", err)
		response.SetError(w, r, err)
		return
	}
	userIDFromJwt, err := request.ReadUserIDFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = validation.ValidateUsersIDs(userID, userIDFromJwt)
	if err != nil {
		u.loggerService.Error("you are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIDToken": userIDFromJwt,
		})
		response.SetError(w, r, err)
		return
	}

	err = u.userService.DeleteUser(r.Context(), userID, userBody.Password)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]string{})
}

func (u *UserController) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	userBody, err := request.ReadBody[DTO.ChangeUserPassword](r)
	if err != nil {
		u.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	userIDString, err := request.ReadParam(r, "userID")
	if err != nil {
		u.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		u.loggerService.Error("failed to covnert string to int", err)
		response.SetError(w, r, err)
		return
	}
	userIDFromJwt, err := request.ReadUserIDFromToken(r)
	if err != nil {
		u.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = validation.ValidateUsersIDs(userID, userIDFromJwt)
	if err != nil {
		u.loggerService.Error("you are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIDToken": userIDFromJwt,
		})
		response.SetError(w, r, err)
		return
	}

	if userBody.NewPassword != userBody.ConfirmPassword {
		err = fmt.Errorf("passwords do not match")
		response.SetError(w, r, err)
		return
	}

	err = u.userService.ChangeUserPassword(r.Context(), userID, userBody.CurrentPassword, userBody.NewPassword)
	if err != nil {
		response.SetError(w, r, err)
	}

	response.Send(w, 204, map[string]string{})
}
