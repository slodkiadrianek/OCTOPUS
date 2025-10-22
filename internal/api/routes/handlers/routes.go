package handlers

import (
	"github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type RouteHandlers struct {
	RouteController *controllers.RouteController
}

func NewRouteHandlers(routeController *controllers.RouteController) *RouteHandlers {
	return &RouteHandlers{
		RouteController: routeController,
	}
}

func (rh *RouteHandlers) SetupRouteHandler(router routes.Router) {
	routeGroup := router.Group("/api/v1/app/:appId/routes")
	//routeGroup.GET("/", rh.RouteController.GetRoutes)
	routeGroup.POST("/", middleware.ValidateMiddleware[DTO.CreateRouteData, *zog.StructSchema]("body", schema.CreateRouteSchema),
		rh.RouteController.AddWorkingRoutes)
	//routeGroup.PUT("/:routeId", rh.RouteController.UpdateRoute)
	//routeGroup.DELETE("/:routeId", rh.RouteController.DeleteRoute)
}
