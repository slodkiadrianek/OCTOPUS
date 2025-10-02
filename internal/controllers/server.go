package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type serverService interface {
	GetServerMetrics(ctx context.Context) ([]models.ServerMetrics, error)
	GetServerInfo() (models.ServerInfo, error)
	InsertServerMetrics(ctx context.Context) error
}
type ServerController struct {
	Logger        *utils.Logger
	ServerService serverService
}

func NewServerController(logger *utils.Logger, serverService serverService) *ServerController {
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
