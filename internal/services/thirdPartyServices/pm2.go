package thirdPartyServices

import (
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type Pm2Service struct {
	appRepository interfaces.AppRepository
	loggerService utils.LoggerService
}

func NewPm2Service(appRepository interfaces.AppRepository, loggerService utils.LoggerService) *Pm2Service {
	return &Pm2Service{
		appRepository: appRepository,
		loggerService: loggerService,
	}
}

// func (pm *Pm2Service) ImportApps() error {
//
// }
