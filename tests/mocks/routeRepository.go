package mocks

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRouteRepository struct {
	mock.Mock
}

func (m *MockRouteRepository) UpdateWorkingRoutesStatuses(ctx context.Context, routeStatuses map[int]string) error {
	args := m.Called(ctx, routeStatuses)
	return args.Error(0)
}
func (m *MockRouteRepository) GetWorkingRoutesToTest(ctx context.Context) ([]models.RouteToTest, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.RouteToTest), args.Error(1)
}

func (m *MockRouteRepository) InsertRoutesInfo(ctx context.Context, routesInfo []*DTO.RouteInfo) ([]int, error) {
	args := m.Called(ctx, routesInfo)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockRouteRepository) InsertRoutesRequests(ctx context.Context, routesRequests []*DTO.RouteRequest) ([]int, error) {
	args := m.Called(ctx, routesRequests)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockRouteRepository) InsertRoutesResponses(ctx context.Context, routesResponses []*DTO.RouteResponse) ([]int, error) {
	args := m.Called(ctx, routesResponses)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockRouteRepository) InsertNextRoutesData(ctx context.Context, nextRoutes []*DTO.NextRoute) ([]int, error) {
	args := m.Called(ctx, nextRoutes)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockRouteRepository) InsertWorkingRoute(ctx context.Context, workingRoute DTO.WorkingRoute) (int, error) {
	args := m.Called(ctx, workingRoute)
	return args.Get(0).(int), args.Error(1)
}
