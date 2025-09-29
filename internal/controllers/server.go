package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type ServerController struct {
	Logger        *utils.Logger
	ServerService *services.ServerService
}

func NewServerController(logger *utils.Logger, serverService *services.ServerService) *ServerController {
	return &ServerController{
		ServerService: serverService,
		Logger:        logger,
	}
}

func (s *ServerController) GetServerInfo(w http.ResponseWriter, r *http.Request) {
	serverInfo, err := s.ServerService.GetServerInfo()
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 200, serverInfo)
}

func (s *ServerController) GetServerMetrics(w http.ResponseWriter, r *http.Request) {
	serverMetrics, err := s.ServerService.GetServerMetrics(r.Context())
	if err != nil {
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
	}
	utils.SendResponse(w, 200, serverMetrics)
}
