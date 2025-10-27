package services

import (
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type Pm2Service struct {
	AppRepository appRepository
	LoggerService *utils.Logger
}

func NewPm2Service(appRepository appRepository, loggerService *utils.Logger) *Pm2Service {
	return &Pm2Service{
		AppRepository: appRepository,
		LoggerService: loggerService,
	}
}

//func (pm *Pm2Service) ImportApps() error {
//
//}
