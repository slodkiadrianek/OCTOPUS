package interfaces

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
)

type RouteRepository interface {
	CheckRouteStatus(ctx context.Context, routeID int) (string, error)
	UpdateWorkingRoutesStatuses(ctx context.Context, routesStatuses map[int]string) error
	GetWorkingRoutesToTest(ctx context.Context) ([]models.RouteToTest, error)
	InsertRoutesInfo(ctx context.Context, routesInfo []*DTO.RouteInfo) ([]int, error)
	InsertRoutesRequests(ctx context.Context,
		routesRequests []*DTO.RouteRequest) ([]int, error)
	InsertRoutesResponses(ctx context.Context,
		routesResponses []*DTO.RouteResponse) ([]int,
		error)
	InsertNextRoutesData(ctx context.Context,
		nextRoutesData []*DTO.NextRoute) ([]int,
		error)
	InsertWorkingRoute(ctx context.Context, workingRoute DTO.WorkingRoute) (int,
		error)
}

type RouteStatusService interface {
	CheckRoutesStatus(ctx context.Context) error
}
