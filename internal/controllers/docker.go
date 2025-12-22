package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
)

type dockerService interface {
	PauseContainer(ctx context.Context, appID string) error
	RestartContainer(ctx context.Context, appID string) error
	StartContainer(ctx context.Context, appID string) error
	UnpauseContainer(ctx context.Context, appID string) error
	StopContainer(ctx context.Context, appID string) error
	ImportContainers(ctx context.Context, ownerID int) error
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
	appID, err := request.ReadParam(r, "appID")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.PauseContainer(r.Context(), appID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]any{})
}

func (dc *DockerController) RestartContainer(w http.ResponseWriter, r *http.Request) {
	appID, err := request.ReadParam(r, "appID")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.RestartContainer(r.Context(), appID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]any{})
}

func (dc *DockerController) StartContainer(w http.ResponseWriter, r *http.Request) {
	appID, err := request.ReadParam(r, "appID")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.StartContainer(r.Context(), appID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]any{})
}

func (dc *DockerController) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	appID, err := request.ReadParam(r, "appID")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.UnpauseContainer(r.Context(), appID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]any{})
}

func (dc *DockerController) StopContainer(w http.ResponseWriter, r *http.Request) {
	appID, err := request.ReadParam(r, "appID")
	if err != nil {
		dc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.StopContainer(r.Context(), appID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 204, map[string]any{})
}

func (dc *DockerController) ImportDockerContainers(w http.ResponseWriter, r *http.Request) {
	ownerID, err := request.ReadUserIDFromToken(r)
	if err != nil {
		dc.loggerService.Error(failedToReadDataFromToken)
		response.SetError(w, r, err)
		return
	}

	err = dc.dockerService.ImportContainers(r.Context(), ownerID)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 201, map[string]any{})
}
