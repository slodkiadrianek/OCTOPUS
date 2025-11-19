package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

type dockerService interface {
	PauseContainer(ctx context.Context, appId string) error
	RestartContainer(ctx context.Context, appId string) error
	StartContainer(ctx context.Context, appId string) error
	UnpauseContainer(ctx context.Context, appId string) error
	StopContainer(ctx context.Context, appId string) error
	ImportContainers(ctx context.Context, ownerId int) error
}

type DockerController struct {
	dockerService dockerService
	loggerService utils.LoggerService
}

func NewDockerController(service dockerService, loggerService utils.LoggerService) *DockerController {
	return &DockerController{
		dockerService: service,
		loggerService: loggerService,
	}
}

func (dc *DockerController) PauseContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	err = dc.dockerService.PauseContainer(r.Context(), appId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) RestartContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	err = dc.dockerService.RestartContainer(r.Context(), appId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) StartContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	err = dc.dockerService.StartContainer(r.Context(), appId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	err = dc.dockerService.UnpauseContainer(r.Context(), appId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]any{})
}
func (dc *DockerController) StopContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		utils.SetError(w, r, err)
		return
	}

	err = dc.dockerService.StopContainer(r.Context(), appId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) ImportDockerContainers(w http.ResponseWriter, r *http.Request) {
	ownerId, err := utils.ReadUserIdFromToken(r)
	if err != nil {
		dc.loggerService.Error(failedToReadDataFromToken)
		utils.SetError(w, r, err)
		return
	}

	err = dc.dockerService.ImportContainers(r.Context(), ownerId)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}

	utils.SendResponse(w, 201, map[string]any{})
}
