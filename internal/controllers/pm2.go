package controllers

import (
	"github.com/slodkiadrianek/octopus/internal/services/thirdPartyServices"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type Pm2Controller struct {
	loggerService utils.LoggerService
	pm2Service    *thirdPartyServices.Pm2Service
}

func NewPm2Controller(loggerService utils.LoggerService, pm2Service *thirdPartyServices.Pm2Service) *Pm2Controller {
	return &Pm2Controller{
		loggerService: loggerService,
		pm2Service:    pm2Service,
	}
}
