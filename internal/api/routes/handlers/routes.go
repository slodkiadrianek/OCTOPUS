package handlers

import (
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/controllers"
)

type RouteHandlers struct {
	RouteController *controllers.RouteController
}

func NewRouteHandlers(routeController *controllers.RouteController) *RouteHandlers {
	return &RouteHandlers{
		RouteController: routeController,
	}
}

func (rh *RouteHandlers) SetupRouteHandler(router *routes.Router) {
	routeGroup := router.Group("/api/v1/app/:appId/routes")
	routeGroup.GET("/", rh.RouteController.GetRoutes)
	routeGroup.POST("/", rh.RouteController.CreateRoute)
	routeGroup.PUT("/:routeId", rh.RouteController.UpdateRoute)
	routeGroup.DELETE("/:routeId", rh.RouteController.DeleteRoute)
}
