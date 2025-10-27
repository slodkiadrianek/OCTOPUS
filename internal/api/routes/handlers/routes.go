package handlers

import (
	"net/http"

	"github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/schema"
)

type routeController interface {
	AddWorkingRoutes(w http.ResponseWriter, r *http.Request)
}
type RouteHandlers struct {
	RouteController routeController
}

func NewRouteHandlers(routeController routeController) *RouteHandlers {
	return &RouteHandlers{
		RouteController: routeController,
	}
}

func (rh *RouteHandlers) SetupRouteHandler(router routes.Router) {
	routeGroup := router.Group("/api/v1/apps/:appId/routes")
	//routeGroup.GET("/", rh.RouteController.GetRoutes)
	routeGroup.POST("/", middleware.ValidateMiddleware[DTO.CreateRouteData, *zog.StructSchema]("body", schema.CreateRouteSchema),
		rh.RouteController.AddWorkingRoutes)
	//routeGroup.PUT("/:routeId", rh.RouteController.UpdateRoute)
	//routeGroup.DELETE("/:routeId", rh.RouteController.DeleteRoute)
}
