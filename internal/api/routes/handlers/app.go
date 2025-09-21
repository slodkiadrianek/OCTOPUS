package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type AppSettingsHandlers struct {
	AppController    *controllers.AppController
	DockerController *controllers.DockerController
	JWT              *middleware.JWT
}

func NewAppAppHandler(appController *controllers.AppController, dockerController *controllers.DockerController, jwt *middleware.JWT) *AppSettingsHandlers {
	return &AppSettingsHandlers{
		AppController:    appController,
		DockerController: dockerController,
		JWT:              jwt,
	}
}

func (a AppSettingsHandlers) SetupAppHandlers(router routes.Router) {
	appGroup := router.Group("/api/v1/app")
	appGroup.POST("", a.JWT.VerifyToken, middleware.ValidateMiddleware[schema.CreateApp]("body", schema.CreateAppSchema), a.AppController.CreateApp)
	appGroup.POST("/docker/import", a.JWT.VerifyToken, a.DockerController.ImportDockerContainers)
	appGroup.GET("/:appId/status", a.JWT.VerifyToken, middleware.ValidateMiddleware[schema.AppId]("params", schema.AppIdSchema), a.AppController.GetAppStatus)
	appGroup.PUT("/:appId")
	appGroup.DELETE("/:appId")
}
