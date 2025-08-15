package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController{
	return &UserController{
		UserService: userService,
	}
}

func (u *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.CreateUser](r)
	if err != nil {
		return
	}
	userDto := DTO.NewUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.UserService.InsertUserToDb(r.Context(), *userDto, userBody.Password)
	if err != nil {
		ctx := context.WithValue(r.Context(), "Error", err)
		r = r.WithContext(ctx)
		return
	}
	utils.SendResponse(w, 201, map[string]string{})
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := utils.ReadBody[schema.UpdateUser](r)
	userId :=1;
	if err != nil{
		return
	}
	userDto := DTO.NewUser(userBody.Email, userBody.Name, userBody.Surname)
	err = u.UserService.UpdateUser(r.Context(), *userDto, userId)

}
