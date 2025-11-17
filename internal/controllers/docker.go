package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type DockerService interface {
	PauseContainer(ctx context.Context, appId string) error
	RestartContainer(ctx context.Context, appId string) error
	StartContainer(ctx context.Context, appId string) error
	UnpauseContainer(ctx context.Context, appId string) error
	StopContainer(ctx context.Context, appId string) error
	ImportContainers(ctx context.Context, ownerId int) error
}

type DockerController struct {
	dockerService DockerService
	loggerService utils.LoggerService
}

func NewDockerController(service DockerService, loggerService utils.LoggerService) *DockerController {
	return &DockerController{
		dockerService: service,
		loggerService: loggerService,
	}
}

func (dc *DockerController) PauseContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	err = dc.dockerService.PauseContainer(r.Context(), appId)
	if err != nil {
		dc.loggerService.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) RestartContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	err = dc.dockerService.RestartContainer(r.Context(), appId)
	if err != nil {
		dc.loggerService.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) StartContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	err = dc.dockerService.StartContainer(r.Context(), appId)
	if err != nil {
		dc.loggerService.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	err = dc.dockerService.UnpauseContainer(r.Context(), appId)
	if err != nil {
		dc.loggerService.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]any{})
}
func (dc *DockerController) StopContainer(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		dc.loggerService.Error("Failed to read app id from params", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
		return
	}
	err = dc.dockerService.StopContainer(r.Context(), appId)
	if err != nil {
		dc.loggerService.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 204, map[string]any{})
}

func (dc *DockerController) ImportDockerContainers(w http.ResponseWriter, r *http.Request) {
	ownerId, ok := r.Context().Value("id").(int)
	if !ok || ownerId == 0 {
		dc.loggerService.Error("Failed to read user id from context", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
	}
	err := dc.dockerService.ImportContainers(r.Context(), ownerId)
	if err != nil {
		dc.loggerService.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 201, map[string]any{})
}
