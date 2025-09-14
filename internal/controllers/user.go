package controllers

import (
	// "context"

	"fmt"
	"net/http"
	"strconv"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (u *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.CreateUser](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	userDto := DTO.NewCreateUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.UserService.InsertUserToDb(r.Context(), *userDto, userBody.Password)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 201, map[string]string{})
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.UpdateUser](r)
	userId := 1
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	userDto := DTO.NewCreateUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.UserService.UpdateUser(r.Context(), *userDto, userId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}
func (u *UserController) UpdateUserNotifications(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.UpdateUserNotifications](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	fmt.Println(1)
	userIdString, err := utils.ReadParam(r, "userId")
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	err = u.UserService.UpdateUserNotifications(r.Context(), userId, *userBody)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.DeleteUser](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	err = u.UserService.DeleteUser(r.Context(), userBody.UserId.UserId, userBody.Password)
	if err != nil {

		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.DeleteUser](r)
	if err != nil {
		return
	}
	err = u.UserService.DeleteUser(r.Context(), userBody.UserId.UserId, userBody.Password)
	if err != nil {
		errBucket ,ok := r.Context().Value("ErrorBucket").(*middleware.ErrorBucket)
		if ok {	
			errBucket.Err = err
	}
		utils.SendResponse(w, 500, map[string]string{
			"errorCategory":    "Server",
			"errorDescription": "Internal server error",
		})
		return
	}
	utils.SendResponse(w, 204, map[string]string{})

}