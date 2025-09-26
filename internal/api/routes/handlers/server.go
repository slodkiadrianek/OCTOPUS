package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type ServerHandlers struct {
	ServerController *controllers.ServerController
	JWT              *middleware.JWT
}

func NewServerHandlers(serverController *controllers.ServerController, jwt *middleware.JWT) *ServerHandlers {
	return &ServerHandlers{
		JWT:              jwt,
		ServerController: serverController,
	}
}

func (s ServerHandlers) SetupServerHandlers(router routes.Router) {
	serverGroup := router.Group("/api/v1/server")
	serverGroup.GET("/metrics", s.JWT.VerifyToken, s.ServerController.GetServerMetrics)
}
