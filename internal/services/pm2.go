package services

import (
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type Pm2Service struct {
	AppRepository appRepository
	Logger        *utils.Logger
}

func NewPm2Service(appRepository appRepository, logger *utils.Logger) *DockerService {
	return &DockerService{
		AppRepository: appRepository,
		Logger:        logger,
	}
}

//func (pm *Pm2Service) ImportApps() error {
//
//}
