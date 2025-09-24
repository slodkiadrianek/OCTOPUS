package controllers

import (
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
	"net/http"
)

type ServerController struct {
	Logger        *logger.Logger
	ServerService *services.ServerService
}

func NewServerController(logger *logger.Logger, serverService *services.ServerService) *ServerController {
	return &ServerController{
		ServerService: serverService,
		Logger:        logger,
	}
}

func (s *ServerController) GetServerMetrics(w http.ResponseWriter, r *http.Request) {
	serverMetrics, err := s.ServerService.GetServerMetrics(r.Context())
	if err != nil {
		err := models.NewError(500, "Server", "Internal server error")
		utils.SetError(w, r, err)
	}
	utils.SendResponse(w, 200, serverMetrics)
}
