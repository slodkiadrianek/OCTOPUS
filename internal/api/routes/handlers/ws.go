package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type WebSocketHandlers struct {
	WsController *controllers.WsController
	JWT          *middleware.JWT
}

func NewWebsocketHandler(wsController *controllers.WsController, jwt *middleware.JWT) *WebSocketHandlers {
	return &WebSocketHandlers{
		WsController: wsController,
		JWT:          jwt,
	}
}

func (ws *WebSocketHandlers) SetupWebsocketHandlers(router routes.Router) {
	groupRouter := router.Group("/ws/v1/apps")
	groupRouter.GET("/:appId/logs", ws.WsController.Logs)
	groupRouter.GET("/:appId/console", ws.JWT.VerifyToken)
}
