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
	jwt *middleware.JWT) *AppSettingsHandlers {
	return &AppSettingsHandlers{
		appController:    appController,
		dockerController: dockerController,
		jwt:              jwt,
	}
}

func (a AppSettingsHandlers) SetupAppHandlers(router routes.Router) {
	appGroup := router.Group("/api/v1/apps")

	appGroup.GET("", a.jwt.VerifyToken, a.appController.GetInfoAboutApps)
	appGroup.GET("/:appId", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params",
		schema.AppIdSchema), a.appController.GetInfoAboutApp)
	appGroup.GET("/:appId/status", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params", schema.AppIdSchema), a.appController.GetAppStatus)

	appGroup.POST("", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.CreateApp]("body", schema.CreateAppSchema),
		a.appController.CreateApp)
	appGroup.POST("/docker/import", a.jwt.VerifyToken, a.dockerController.ImportDockerContainers)

	appGroup.PUT("/:appId/docker/stop", a.jwt.VerifyToken, a.dockerController.StopContainer)
	appGroup.PUT("/:appId/docker/start", a.jwt.VerifyToken, a.dockerController.StartContainer)
	appGroup.PUT("/:appId/docker/restart", a.jwt.VerifyToken, a.dockerController.RestartContainer)
	appGroup.PUT("/:appId/docker/pause", a.jwt.VerifyToken, a.dockerController.PauseContainer)
	appGroup.PUT("/:appId/docker/unpause", a.jwt.VerifyToken, a.dockerController.UnpauseContainer)

	appGroup.PUT("/:appId", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params", schema.AppIdSchema),
		middleware.ValidateMiddleware[DTO.UpdateApp]("body", schema.UpdateAppSchema), a.appController.UpdateApp)
	appGroup.DELETE("/:appId", a.jwt.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params",
		schema.AppIdSchema), a.appController.DeleteApp)

}
