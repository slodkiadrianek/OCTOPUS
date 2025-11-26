package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
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
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.PauseContainer(r.Context(), appId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) RestartContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.RestartContainer(r.Context(), appId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) StartContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.StartContainer(r.Context(), appId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.UnpauseContainer(r.Context(), appId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.SendResponse(w, 204, map[string]any{})
}
func (dc *DockerController) StopContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.StopContainer(r.Context(), appId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) ImportDockerContainers(w http.ResponseWriter, r *http.Request) {
	ownerId, err := request.ReadUserIdFromToken(r)
	if err != nil {
		dc.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.ImportContainers(r.Context(), ownerId)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.SendResponse(w, 201, map[string]any{})
}
