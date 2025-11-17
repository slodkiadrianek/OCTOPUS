package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type appService interface {
	CreateApp(ctx context.Context, app DTO.CreateApp, ownerId int) error
	GetApp(ctx context.Context, id string, ownerId int) (*models.App, error)
	GetApps(ctx context.Context, ownerId int) ([]models.App, error)
	GetAppStatus(ctx context.Context, id string, ownerId int) (DTO.AppStatus, error)
	DeleteApp(ctx context.Context, id string, ownerId int) error
	CheckAppsStatus(ctx context.Context) ([]DTO.AppStatus, error)
	SendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) error
	UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error
}

type AppController struct {
	appService    appService
	loggerService utils.LoggerService
}

func NewAppController(appService appService, loggerService utils.LoggerService) *AppController {
	return &AppController{
		appService:    appService,
		loggerService: loggerService,
	}
}

func (a *AppController) GetInfoAboutApps(w http.ResponseWriter, r *http.Request) {
	ownerId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if ownerId == 0 {
		return
	}
	apps, err := a.appService.GetApps(r.Context(), ownerId)
	if err != nil {
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
	ownerId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if ownerId == 0 {
		return
	}
	app, err := a.appService.GetApp(r.Context(), appId, ownerId)
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
	ownerId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if ownerId == 0 {
		return
	}
	err = a.appService.CreateApp(r.Context(), *appBody, ownerId)
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
	ownerId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if ownerId == 0 {
		return
	}
	err = a.appService.UpdateApp(r.Context(), appId, *app, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}

func (a *AppController) DeleteApp(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		a.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if ownerId == 0 {
		return
	}
	err = a.appService.DeleteApp(r.Context(), appId, ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]string{})
}

func (a *AppController) GetAppStatus(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		a.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	ownerId := utils.ReadUserIdFromToken(w, r, a.loggerService)
	if ownerId == 0 {
		return
	}
	appStatus, err := a.appService.GetAppStatus(r.Context(), appId, ownerId)
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
		a.loggerService.Error("Failed to read app id from context", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	// status, err := a.appService.GetDbStatus(r.Context(), appId)
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
// 		a.loggerService.Error("Failed to read app id from context", r.URL.Path)
// 		err := models.NewError(500, "Server", "Internal server error")
// 		utils.SetError(w, r, err)
// 		return
// 	}
// 	metrics , err := a.appService.GetServerMetrics(r.Context(), appId)
// 	if err != nil {
// 		utils.SetError(w, r, err)
// 		return
// 	}
// 	utils.SendResponse(w, 200, map[string]any{
// 		"metrics": metrics,
// 	})
// }
