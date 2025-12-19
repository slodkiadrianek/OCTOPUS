package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
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
	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	apps, err := a.appService.GetApps(r.Context(), ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, apps)
}

func (a *AppController) GetInfoAboutApp(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	app, err := a.appService.GetApp(r.Context(), appId, ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, app)
}

func (a *AppController) CreateApp(w http.ResponseWriter, r *http.Request) {
	appBody, err := request.ReadBody[DTO.CreateApp](r)
	if err != nil {
		a.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = a.appService.CreateApp(r.Context(), *appBody, ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 201, map[string]string{})
}

func (a *AppController) UpdateApp(w http.ResponseWriter, r *http.Request) {
	app, err := request.ReadBody[DTO.UpdateApp](r)
	if err != nil {
		a.loggerService.Error(failedToReadBodyFromRequest, err)
		response.SetError(w, r, err)
		return
	}

	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		a.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = a.appService.UpdateApp(r.Context(), appId, *app, ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]string{})
}

func (a *AppController) DeleteApp(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		a.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = a.appService.DeleteApp(r.Context(), appId, ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]string{})
}

func (a *AppController) GetAppStatus(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		a.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		a.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	appStatus, err := a.appService.GetAppStatus(r.Context(), appId, ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, map[string]DTO.AppStatus{
		"data": appStatus,
	})
}
