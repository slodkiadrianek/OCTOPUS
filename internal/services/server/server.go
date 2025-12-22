package server

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type ServerService struct {
	logger       utils.LoggerService
	cacheService interfaces.CacheService
}

func NewServerService(logger utils.LoggerService, cacheService interfaces.CacheService) *ServerService {
	return &ServerService{
		logger:       logger,
		cacheService: cacheService,
	}
}

func (s *ServerService) GetServerMetrics(ctx context.Context) ([]models.ServerMetrics, error) {
	serverMetricsData, err := s.cacheService.GetData(ctx, "server_metrics")
	if err != nil {
		s.logger.Error("Failed to get server metrics from cache", err)
		return nil, err
	}

	existingMetrics, err := utils.UnmarshalData[[]models.ServerMetrics]([]byte(serverMetricsData))
	if err != nil {
		s.logger.Error("Failed to unmarshal server metrics data", err)
		return nil, err
	}

	return *existingMetrics, nil
}

func (s *ServerService) GetServerInfo() (models.ServerInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		s.logger.Warn("Failed to read cpu data", err)
		return models.ServerInfo{}, err
	}

	ram, err := mem.VirtualMemory()
	if err != nil {
		s.logger.Warn("Failed to read memory data", err)
		return models.ServerInfo{}, err
	}

	diskUsage, err := disk.Usage("./")
	if err != nil {
		s.logger.Warn("Failed to read disk data", err)
		return models.ServerInfo{}, err
	}

	workTime, err := host.Uptime()
	if err != nil {
		s.logger.Warn("Failed to uptime  data", err)
		return models.ServerInfo{}, err
	}

	info, err := host.Info()
	if err != nil {
		s.logger.Warn("Failed to info data", err)
		return models.ServerInfo{}, err
	}

	modelName := cpuInfo[0].ModelName
	uptime := workTime / 60
	ramTotal := ram.Total / 1000 / 1000 / 1000
	diskTotal := diskUsage.Total / 1000 / 1000 / 1000
	serverInfo := models.NewServerInfo(modelName, info.OS, info.Platform, info.KernelVersion, ramTotal, diskTotal, uptime)

	return *serverInfo, nil
}

func (s *ServerService) InsertServerMetrics(ctx context.Context) error {
	cpuData, err := cpu.Percent(time.Second, false)
	if err != nil {
		s.logger.Warn("Failed to read cpu data", err)
		return err
	}

	memData, err := mem.VirtualMemory()
	if err != nil {
		s.logger.Warn("Failed to read memory data", err)
		return err
	}

	memPercent := memData.UsedPercent
	diskUsage, err := disk.Usage("./")
	if err != nil {
		s.logger.Warn("Failed to read disk data", err)
		return err
	}

	actualDate := time.Now()
	serverMetricsData := models.NewServerMetrics(int(cpuData[0]), int(memPercent), int(diskUsage.UsedPercent), actualDate)
	s.logger.Info("Server metrics data", map[string]any{
		"CPU":  int(cpuData[0]),
		"RAM":  int(memPercent),
		"Disk": int(diskUsage.UsedPercent),
		"Date": actualDate,
	})
	doesServerMetricsExists, err := s.cacheService.ExistsData(ctx, "server_metrics")
	var existingMetrics *[]models.ServerMetrics
	if err != nil {
		s.logger.Warn("Failed to check if server metrics exist in cache", err)
		return err
	}
	if doesServerMetricsExists == 1 {
		serverMetricsData, err := s.cacheService.GetData(ctx, "server_metrics")
		if err != nil {
			s.logger.Warn("Failed to get server metrics from cache", err)
			return err
		}
		existingMetrics, err = utils.UnmarshalData[[]models.ServerMetrics]([]byte(serverMetricsData))
		if err != nil {
			return err
		}
	}

	if existingMetrics != nil {
		if len(*existingMetrics) >= 40 {
			*existingMetrics = (*existingMetrics)[1:]
		}
		*existingMetrics = append(*existingMetrics, *serverMetricsData)
	} else {
		existingMetrics = &[]models.ServerMetrics{*serverMetricsData}
	}

	bodyBytes, err := utils.MarshalData(existingMetrics)
	if err != nil {
		s.logger.Warn("Failed to marshal server metrics data", err)
		return err
	}

	err = s.cacheService.SetData(ctx, "server_metrics", string(bodyBytes), 0)
	if err != nil {
		s.logger.Warn("Failed to cache server metrics", err)
		return err
	}

	return nil
}
