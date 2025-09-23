package services

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
	"time"
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

func (s *ServerService) GetServerMetrics() ([]models.ServerMetrics, error) {

}

func (s *ServerService) InsertServerMetrics() error {
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
	diskUsage, err := disk.Usage("/")
	if err != nil {
		s.Logger.Warn("Failed to read disk data", err)
		return err
	}
}
