package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/DTO"
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
	appGroup := router.Group("/api/v1/apps")
	appGroup.POST("", a.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.CreateApp]("body", schema.CreateAppSchema),
		a.AppController.CreateApp)
	appGroup.POST("/docker/import", a.JWT.VerifyToken, a.DockerController.ImportDockerContainers)
	//appGroup.GET("/:appId/logs", a.JWT.VerifyToken, a.AppController.GetLogs)
	appGroup.GET("/:appId/status", a.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params",
		schema.AppIdSchema), a.AppController.GetAppStatus)
	appGroup.PUT("/:appId", a.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params", schema.AppIdSchema),
		middleware.ValidateMiddleware[DTO.UpdateApp]("body", schema.UpdateAppSchema), a.AppController.UpdateApp)
	appGroup.DELETE("/:appId", a.JWT.VerifyToken, middleware.ValidateMiddleware[DTO.AppId]("params",
		schema.AppIdSchema), a.AppController.DeleteApp)

}
