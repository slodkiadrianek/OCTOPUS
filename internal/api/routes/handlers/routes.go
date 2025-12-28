package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/interfaces"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type RouteHandlers struct {
	routeController interfaces.RouteController
}

func NewRouteHandlers(routeController interfaces.RouteController) *RouteHandlers {
	return &RouteHandlers{
		routeController: routeController,
	}
}

func (rh *RouteHandlers) SetupRouteHandler(router routes.Router) {
	routeGroup := router.Group("/api/v1/apps/:appID/routes")

	routeGroup.GET("/:routeID", middleware.ValidateMiddleware[DTO.RouteID]("params", schema.RouteIDSchema), rh.routeController.CheckRouteStatus)
	routeGroup.POST("/", middleware.ValidateMiddleware[DTO.CreateRouteData]("body", schema.CreateRouteSchema),
		rh.routeController.AddWorkingRoutes)
	// routeGroup.PUT("/:routeId", rh.routeController.UpdateRoute)
	// routeGroup.DELETE("/:routeId", rh.routeController.DeleteRoute)
}
