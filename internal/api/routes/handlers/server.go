package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type ServerHandlers struct {
	serverController interfaces.ServerController
	jwt              *middleware.JWT
}

func NewServerHandlers(serverController interfaces.ServerController, jwt *middleware.JWT) *ServerHandlers {
	return &ServerHandlers{
		jwt:              jwt,
		serverController: serverController,
	}
}

func (s ServerHandlers) SetupServerHandlers(router routes.Router) {
	serverGroup := router.Group("/api/v1/server")

	serverGroup.GET("", s.jwt.VerifyToken, s.serverController.GetServerInfo)
	serverGroup.GET("/metrics", s.jwt.VerifyToken, s.serverController.GetServerMetrics)
}
