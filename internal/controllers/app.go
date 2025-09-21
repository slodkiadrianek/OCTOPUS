package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type key string

type AppController struct {
	AppService *services.AppService
	Logger     *logger.Logger
}

func NewAppController(appService *services.AppService, logger *logger.Logger) *AppController {
	return &AppController{
		AppService: appService,
		Logger:     logger,
	}
}

func (a *AppController) CreateApp(w http.ResponseWriter, r *http.Request) {
	appBody, err := utils.ReadBody[schema.CreateApp](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	ownerId, ok := r.Context().Value("id").(int)
	if !ok || ownerId == 0 {
		a.Logger.Error("Failed to read user id from context", r.URL.Path)
		err = models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
	}
	err = a.AppService.CreateApp(r.Context(), *appBody, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 201, map[string]string{})
}

func (a *AppController) GetApp(w http.ResponseWriter, r *http.Request) {}

func (a *AppController) UpdateApp(w http.ResponseWriter, r *http.Request) {}

func (a *AppController) DeleteApp(w http.ResponseWriter, r *http.Request) {}

func (a *AppController) GetAppStatus(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		a.Logger.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	appStatus, err := a.AppService.GetAppStatus(r.Context(), appId)
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
