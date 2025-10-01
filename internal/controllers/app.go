package controllers

import (
	"fmt"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AppController struct {
	AppService *services.AppService
	Logger     *utils.Logger
}

func NewAppController(appService *services.AppService, logger *utils.Logger) *AppController {
	return &AppController{
		AppService: appService,
		Logger:     logger,
	}
}

func (a *AppController) GetInfoAboutApps(w http.ResponseWriter, r *http.Request) {
	ownerId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if ownerId == 0 {
		return
	}
	apps, err := a.AppService.GetApps(r.Context(), ownerId)
	fmt.Println("TEST")
	if err != nil {
		fmt.Println(err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 200, apps)
}

func (a *AppController) GetInfoAboutApp(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if ownerId == 0 {
		return
	}
	app, err := a.AppService.GetApp(r.Context(), appId, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 200, app)
}

func (a *AppController) CreateApp(w http.ResponseWriter, r *http.Request) {
	appBody, err := utils.ReadBody[DTO.CreateApp](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if ownerId == 0 {
		return
	}
	err = a.AppService.CreateApp(r.Context(), *appBody, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 201, map[string]string{})
}

func (a *AppController) UpdateApp(w http.ResponseWriter, r *http.Request) {
	app, err := utils.ReadBody[DTO.UpdateApp](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if ownerId == 0 {
		return
	}
	err = a.AppService.UpdateApp(r.Context(), appId, *app, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}

func (a *AppController) DeleteApp(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		a.Logger.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if ownerId == 0 {
		return
	}
	err = a.AppService.DeleteApp(r.Context(), appId, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}

func (a *AppController) GetAppStatus(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		a.Logger.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.Logger)
	if ownerId == 0 {
		return
	}
	appStatus, err := a.AppService.GetAppStatus(r.Context(), appId, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 200, map[string]DTO.AppStatus{
		"data": appStatus,
	})
}

func (a *AppController) GetDbStatus(w http.ResponseWriter, r *http.Request) {
	appId, ok := r.Context().Value("appId").(int)
	if !ok || appId == 0 {
		a.Logger.Error("Failed to read app id from context", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	// status, err := a.AppService.GetDbStatus(r.Context(), appId)
	// if err != nil {
	// 	utils.SetError(w, r, err)
	// 	return
	// }
	utils.SendResponse(w, 200, map[string]string{
		// "status": status,
	})
}

// func (a *AppController) GetServerMetrics(w http.ResponseWriter, r *http.Request){
// 	appId,  err := utils.ReadParam(r, "appId")
// 	if err != nil {
// 		a.Logger.Error("Failed to read app id from context", r.URL.Path)
// 		err := models.NewError(500, "Server", "Internal server error")
// 		utils.SetError(w, r, err)
// 		return
// 	}
// 	metrics , err := a.AppService.GetServerMetrics(r.Context(), appId)
// 	if err != nil {
// 		utils.SetError(w, r, err)
// 		return
// 	}
// 	utils.SendResponse(w, 200, map[string]any{
// 		"metrics": metrics,
// 	})
// }
