package controllers

import (
	"github.com/slodkiadrianek/octopus/internal/services/thirdPartyServices"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type Pm2Controller struct {
	Logger     *utils.Logger
	Pm2Service *thirdPartyServices.Pm2Service
}

func NewPm2Controller(logger *utils.Logger, pm2Service *thirdPartyServices.Pm2Service) *Pm2Controller {
	return &Pm2Controller{
		Logger:     logger,
		Pm2Service: pm2Service,
	}
}
