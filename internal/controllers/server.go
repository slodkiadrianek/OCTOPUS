package controllers

import (
	"github.com/slodkiadrianek/octopus/internal/services"
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

}
