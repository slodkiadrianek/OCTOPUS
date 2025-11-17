package handlers

import (
	"github.com/Oudwins/zog"
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
	routeGroup := router.Group("/api/v1/apps/:appId/routes")
	//routeGroup.GET("/", rh.routeController.GetRoutes)
	routeGroup.POST("/", middleware.ValidateMiddleware[DTO.CreateRouteData, *zog.StructSchema]("body", schema.CreateRouteSchema),
		rh.routeController.AddWorkingRoutes)
	//routeGroup.PUT("/:routeId", rh.routeController.UpdateRoute)
	//routeGroup.DELETE("/:routeId", rh.routeController.DeleteRoute)
}
