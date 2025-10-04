package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type wsService interface {
	Logs(ctx context.Context, appId string, conn *websocket.Conn)
}

type WsController struct {
	WsService wsService
	Logger    *utils.Logger
}

func NewWsController(wsService wsService, logger *utils.Logger) *WsController {
	return &WsController{
		WsService: wsService,
		Logger:    logger,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WsController) Logs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.Logger.Error("Failed to upgrade connection", err)
		return
	}
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		ws.Logger.Error("Failed to read param", appId)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to read appId"))
		return
	}
	ctx := context.Background()
	ws.WsService.Logs(ctx, appId, conn)
}
