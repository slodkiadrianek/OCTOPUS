package services

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type ServerService struct {
	Logger       *logger.Logger
	CacheService *config.CacheService
}

func NewServerService(logger *logger.Logger, cacheService *config.CacheService) *ServerService {
	return &ServerService{
		Logger:       logger,
		CacheService: cacheService,
	}
}

func (s *ServerService) GetServerMetrics(ctx context.Context) ([]models.ServerMetrics, error) {
	serverMetricsData, err := s.CacheService.GetData(ctx, "server_metrics")
	if err != nil {
		s.Logger.Warn("Failed to get server metrics from cache", err)
		return nil, err
	}
	existingMetrics, err := utils.UnmarshalData[[]models.ServerMetrics]([]byte(serverMetricsData))
	if err != nil {
		s.Logger.Warn("Failed to unmarshal server metrics data", err)
		return nil, err
	}
	return *existingMetrics, nil
}

func (s *ServerService) InsertServerMetrics(ctx context.Context) error {
	cpuData, err := cpu.Percent(time.Second, false)
	if err != nil {
		s.Logger.Warn("Failed to read cpu data", err)
		return err
	}
	memData, err := mem.VirtualMemory()
	if err != nil {
		s.Logger.Warn("Failed to read memory data", err)
		return err
	}
	memPercent := memData.UsedPercent
	diskUsage, err := disk.Usage("./")
	if err != nil {
		s.Logger.Warn("Failed to read disk data", err)
		return err
	}
	actualDate := time.Now()
	metricsData := models.NewServerMetrics(int(cpuData[0]), int(memPercent), int(diskUsage.UsedPercent), actualDate)
	s.Logger.Info("Server metrics data", map[string]any{
		"CPU":  int(cpuData[0]),
		"RAM":  int(memPercent),
		"Disk": int(diskUsage.UsedPercent),
		"Date": actualDate,
	})
	//fmt.Println(metricsData)
	doesExist, err := s.CacheService.ExistsData(ctx, "server_metrics")
	var existingMetrics *[]models.ServerMetrics
	if err != nil {
		s.Logger.Warn("Failed to check if server metrics exist in cache", err)
		return err
	}
	if doesExist == 1 {
		serverMetricsData, err := s.CacheService.GetData(ctx, "server_metrics")
		if err != nil {
			s.Logger.Warn("Failed to get server metrics from cache", err)
			return err
		}
		existingMetrics, err = utils.UnmarshalData[[]models.ServerMetrics]([]byte(serverMetricsData))
	}
	if existingMetrics != nil {
		if len(*existingMetrics) >= 40 {
			*existingMetrics = (*existingMetrics)[1:]
		}
		*existingMetrics = append(*existingMetrics, *metricsData)
	} else {
		existingMetrics = &[]models.ServerMetrics{*metricsData}
	}
	bodyBytes, err := utils.MarshalData(existingMetrics)
	err = s.CacheService.SetData(ctx, "server_metrics", string(bodyBytes), 0)
	if err != nil {
		s.Logger.Warn("Failed to cache server metrics", err)
		return err
	}
	return nil
}
