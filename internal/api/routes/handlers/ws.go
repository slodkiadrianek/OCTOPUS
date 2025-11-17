package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type WebSocketHandlers struct {
	wsController interfaces.WsController
	jwt          *middleware.JWT
}

func NewWebsocketHandler(wsController interfaces.WsController, jwt *middleware.JWT) *WebSocketHandlers {
	return &WebSocketHandlers{
		wsController: wsController,
		jwt:          jwt,
	}
}

func (ws *WebSocketHandlers) SetupWebsocketHandlers(router routes.Router) {
	groupRouter := router.Group("/ws/v1/apps")
	groupRouter.GET("/:appId/logs", ws.wsController.Logs)
	//groupRouter.GET("/:appId/console", ws.jwt.VerifyToken)
}
