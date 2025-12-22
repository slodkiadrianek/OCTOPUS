package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type AppSettingsHandlers struct {
	appController    interfaces.AppController
	dockerController interfaces.DockerController
	jwt              *middleware.JWT
}

func NewAppAppHandler(appController interfaces.AppController, dockerController interfaces.DockerController,
	jwt *middleware.JWT,
) *AppSettingsHandlers {
	return &AppSettingsHandlers{
		appController:    appController,
		dockerController: dockerController,
		jwt:              jwt,
	}
}

func (a AppSettingsHandlers) SetupAppHandlers(router routes.Router) {
	appGroup := router.Group("/api/v1/apps")

	appGroup.GET("", a.jwt.VerifyToken, a.appController.GetInfoAboutApps)
	appGroup.GET("/:appID", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppID]("params",
		schema.AppIDSchema), a.appController.GetInfoAboutApp)
	appGroup.GET("/:appID/status", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppID]("params", schema.AppIDSchema), a.appController.GetAppStatus)

	appGroup.POST("", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.CreateApp]("body", schema.CreateAppSchema),
		a.appController.CreateApp)
	appGroup.POST("/docker/import", a.jwt.VerifyToken, a.dockerController.ImportDockerContainers)

	appGroup.PUT("/:appID/docker/stop", a.jwt.VerifyToken, a.dockerController.StopContainer)
	appGroup.PUT("/:appID/docker/start", a.jwt.VerifyToken, a.dockerController.StartContainer)
	appGroup.PUT("/:appID/docker/restart", a.jwt.VerifyToken, a.dockerController.RestartContainer)
	appGroup.PUT("/:appID/docker/pause", a.jwt.VerifyToken, a.dockerController.PauseContainer)
	appGroup.PUT("/:appID/docker/unpause", a.jwt.VerifyToken, a.dockerController.UnpauseContainer)

	appGroup.PUT("/:appID", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppID]("params", schema.AppIDSchema),
		middleware.ValidateMiddleware[DTO.UpdateApp]("body", schema.UpdateAppSchema), a.appController.UpdateApp)
	appGroup.DELETE("/:appID", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppID]("params",
		schema.AppIDSchema), a.appController.DeleteApp)
}
