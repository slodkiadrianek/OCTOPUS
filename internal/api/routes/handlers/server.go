package handlers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type serverController interface {
	GetServerInfo(w http.ResponseWriter, r *http.Request)
	GetServerMetrics(w http.ResponseWriter, r *http.Request)
}
type ServerHandlers struct {
	ServerController serverController
	JWT              *middleware.JWT
}

func NewServerHandlers(serverController serverController, jwt *middleware.JWT) *ServerHandlers {
	return &ServerHandlers{
		JWT:              jwt,
		ServerController: serverController,
	}
}

func (s ServerHandlers) SetupServerHandlers(router routes.Router) {
	serverGroup := router.Group("/api/v1/server")
	serverGroup.GET("", s.JWT.VerifyToken, s.ServerController.GetServerInfo)
	serverGroup.GET("/metrics", s.JWT.VerifyToken, s.ServerController.GetServerMetrics)
}
