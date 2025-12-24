package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
)

type wsService interface {
	Logs(ctx context.Context, appID string, conn *websocket.Conn)
}

type WsController struct {
	wsService     wsService
	loggerService utils.LoggerService
}

func NewWsController(wsService wsService, loggerService utils.LoggerService) *WsController {
	return &WsController{
		wsService:     wsService,
		loggerService: loggerService,
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
		ws.loggerService.Error("failed to upgrade connection", err)
		return
	}

	appID, err := request.ReadParam(r, "appID")
	if err != nil {
		ws.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		err := conn.WriteMessage(websocket.TextMessage, []byte("failed to read appID"))
		if err != nil {
			return
		}

		return
	}

	ctx := context.Background()
	ws.wsService.Logs(ctx, appID, conn)
}
