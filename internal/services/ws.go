package services

import (
	"bufio"
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type WsService struct {
	DockerHost string
	Logger     *utils.Logger
}

func NewWsService(logger *utils.Logger, dockerHost string) *WsService {
	return &WsService{
		DockerHost: dockerHost,
		Logger:     logger,
	}
}

func (ws *WsService) Logs(ctx context.Context, appId string, conn *websocket.Conn) {
	cli, err := client.NewClientWithOpts(client.WithHost(ws.DockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		ws.Logger.Error("Error creating Docker client", err)
		return
	}
	defer cli.Close()

	conn.SetPongHandler(func(string) error {
		return nil
	})

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
		Tail:       "100",
	}

	reader, err := cli.ContainerLogs(ctx, appId, options)
	if err != nil {
		ws.Logger.Error("Failed to connect to container", err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed to connect to container: %v", err)))
		return
	}
	defer reader.Close()

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	scanner := bufio.NewScanner(reader)

	for {
		select {
		case <-pingTicker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				ws.Logger.Info("Client disconnected during ping")
				return
			}
		case <-ctx.Done():
			ws.Logger.Info("Status check stopped by context")
			return
		default:
			if scanner.Scan() {
				line := scanner.Text()
				clean := make([]rune, 0, len(line))
				for _, r := range line {
					if r == utf8.RuneError {
						continue
					}
					clean = append(clean, r)
				}
				utf8Bytes := []byte(string(clean))
				if err := conn.WriteMessage(websocket.TextMessage, utf8Bytes); err != nil {
					if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						ws.Logger.Info("Client disconnected")
					}
					return
				}
			} else {
				if err := scanner.Err(); err != nil {
					ws.Logger.Error("Scanner error", err)
				}
				return
			}
		}
	}
}
