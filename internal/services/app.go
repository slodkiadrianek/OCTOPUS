package services

import (
	"context"
	"net"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type AppService struct {
	AppRepository *repository.AppRepository
	Logger        *logger.Logger
	CacheService  *config.CacheService
}

func NewAppService(appRepository *repository.AppRepository, logger *logger.Logger, cacheService *config.CacheService) *AppService {
	return &AppService{
		AppRepository: appRepository,
		Logger:        logger,
		CacheService:  cacheService,
	}
}

func (a *AppService) CreateApp(ctx context.Context, app schema.CreateApp, ownerId int) error {
	appDto := DTO.NewApp(app.Name, app.Description, app.DbLink, app.ApiLink, ownerId, app.DiscordWebhook, app.SlackWebhook)
	err := a.AppRepository.InsertApp(ctx, *appDto)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetApp(ctx context.Context, id int) (*models.App, error) {
	app, err := a.AppRepository.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *AppService) UpdateApp(ctx context.Context, id int, app schema.UpdateApp) error {
	appDto := DTO.NewUpdateApp(id, app.Name, app.Description, app.DbLink, app.ApiLink, app.DiscordWebhook, app.SlackWebhook)
	err := a.AppRepository.UpdateApp(ctx, *appDto)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) DeleteApp(ctx context.Context, id int) error {
	err := a.AppRepository.DeleteApp(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) GetAppStatus(ctx context.Context, id int) (string, error) {
	appServerAddress, err := a.AppRepository.GetAppServerAddress(ctx, id)
	if err != nil {
		a.Logger.Error("Failed to get app server address from database", err)
		return "", models.NewError(400, "app server", "Failed to get info about app status")
	}
	packetConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil{
		a.Logger.Error("Failed to listen to packet connection", err)
		return "", models.NewError(400, "app server", "Failed to get info about app status");
	}
	defer packetConn.Close()
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte("OCTOPUS"),
		},
	}
	wb, err := msg.Marshal(nil)
	if err != nil {
		a.Logger.Error("Failed to marshal ICMP message", err)
		return "", models.NewError(400, "app server", "Failed to get info about app status")
	}
	if _, err := packetConn.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(appServerAddress)}); err != nil {
		a.Logger.Error("Failed to write to packet connection", err)
		return "", models.NewError(400, "app server", "Failed to get info about app status")
	}
	rb:=make([]byte, 1500)
	n,peer,err := packetConn.ReadFrom(rb)
	if err != nil {
		a.Logger.Error("Failed to read from packet connection", err)
		return "", models.NewError(400, "app server", "Failed to get info about app status")
	}
	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		a.Logger.Error("Failed to parse ICMP message", err)
		return "", models.NewError(400, "app server", "Failed to get info about app status")
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		a.Logger.Info("Got echo reply from " + peer.String())
		return "Online", nil
	default:
		a.Logger.Info("Got something else from " + peer.String())
		return "Offline", nil
	}
}

func (a *AppService) GetDbStatus(ctx context.Context, id int) (string, error) {
	dbServerAddress, err := a.AppRepository.GetDbServerAddress(ctx, id)
	if err != nil {
		a.Logger.Error("Failed to get db server address from database", err)
		return "", models.NewError(400, "db server", "Failed to get info about db status")
	}
	conn, err := net.Dial("tcp", dbServerAddress)
	if err != nil {
		a.Logger.Error("Failed to connect to db server", err)
		return "Offline", nil
	}
	defer conn.Close()
	return "Online", nil
}