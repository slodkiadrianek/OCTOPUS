package handlers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type wsController interface {
	Logs(w http.ResponseWriter, r *http.Request)
}
type WebSocketHandlers struct {
	WsController wsController
	JWT          *middleware.JWT
}

func NewWebsocketHandler(wsController wsController, jwt *middleware.JWT) *WebSocketHandlers {
	return &WebSocketHandlers{
		WsController: wsController,
		JWT:          jwt,
	}
}

func (ws *WebSocketHandlers) SetupWebsocketHandlers(router routes.Router) {
	groupRouter := router.Group("/ws/v1/apps")
	groupRouter.GET("/:appId/logs", ws.WsController.Logs)
	//groupRouter.GET("/:appId/console", ws.JWT.VerifyToken)
}
