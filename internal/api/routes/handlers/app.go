package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
)

type AppSettingsHandlers struct {
	AppController *controllers.AppController
}

func NewAppAppSettingsHandlers(appController *controllers.AppController) *AppSettingsHandlers {
	return &AppSettingsHandlers{
		AppController: appController,
	}
}

func (a AppSettingsHandlers) SetupAppHandlers(router *routes.Router) {
	appGroup := router.Group("/api/v1/app")
	appGroup.POST("")
	appGroup.GET("/:appId")
	appGroup.PUT("/:appId")
	appGroup.DELETE("/:appId")
}

