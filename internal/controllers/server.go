package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
)

type serverService interface {
	GetServerMetrics(ctx context.Context) ([]models.ServerMetrics, error)
	GetServerInfo() (models.ServerInfo, error)
	InsertServerMetrics(ctx context.Context) error
}
type ServerController struct {
	loggerService utils.LoggerService
	serverService serverService
}

func NewServerController(loggerService utils.LoggerService, serverService serverService) *ServerController {
	return &ServerController{
		serverService: serverService,
		loggerService: loggerService,
	}
}

func (s *ServerController) GetServerInfo(w http.ResponseWriter, r *http.Request) {
	serverInfo, err := s.serverService.GetServerInfo()
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, serverInfo)
}

func (s *ServerController) GetServerMetrics(w http.ResponseWriter, r *http.Request) {
	serverMetrics, err := s.serverService.GetServerMetrics(r.Context())
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 200, serverMetrics)
}
