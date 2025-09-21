package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type DockerController struct {
	DockerService *services.DockerService
	Logger        *logger.Logger
}

func NewDockerController(service *services.DockerService, logger *logger.Logger) *DockerController {
	return &DockerController{
		DockerService: service,
		Logger:        logger,
	}
}

func (dc *DockerController) ImportDockerContainers(w http.ResponseWriter, r *http.Request) {
	ownerId, ok := r.Context().Value("id").(int)
	if !ok || ownerId == 0 {
		dc.Logger.Error("Failed to read user id from context", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
	}
	err := dc.DockerService.ImportContainers(r.Context(), ownerId)
	if err != nil {
		dc.Logger.Error("Failed to import docker containers", r.URL.Path, err)
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 201, map[string]any{})
}
