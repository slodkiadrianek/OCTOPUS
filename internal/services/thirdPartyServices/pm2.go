package thirdPartyServices

import (
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type Pm2Service struct {
	AppRepository interfaces.AppRepository
	LoggerService *utils.Logger
}

func NewPm2Service(appRepository interfaces.AppRepository, loggerService *utils.Logger) *Pm2Service {
	return &Pm2Service{
		AppRepository: appRepository,
		LoggerService: loggerService,
	}
}

// func (pm *Pm2Service) ImportApps() error {
//
// }
