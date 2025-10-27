package services

import (
	"context"
	"runtime"
	"sync"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type routeRepository interface {
	UpdateWorkingRoutesStatuses(ctx context.Context, routesStatuses map[int]string) error
	GetWorkingRoutesToTest(ctx context.Context) ([]DTO.RouteToTest, error)
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

type RouteService struct {
	Logger          *utils.Logger
	RouteRepository routeRepository
}

func NewRouteService(logger *utils.Logger, routeRepository routeRepository) *RouteService {
	return &RouteService{
		Logger:          logger,
		RouteRepository: routeRepository,
	}
}

func (rs *RouteService) AddWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appId string, name string) error {
	nextRoutesChan := make(chan DTO.NextRoute, len(*routes))
	requestRoutesChan := make(chan DTO.RouteRequest, len(*routes))
	responseRoutesChan := make(chan DTO.RouteResponse, len(*routes))
	routesInfoChan := make(chan DTO.RouteInfo, len(*routes))
	errorChan := make(chan error)
	jobs := make(chan DTO.CreateRoute, len(*routes))
	workerCount := runtime.NumCPU()
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				responseBodyBytes, err := utils.MarshalData(job.ResponseBody)
				if err != nil {
					errorChan <- err
					return
				}
				responseRoutesChan <- *DTO.NewRouteResponse(job.ResponseStatusCode, job.ParentId, string(responseBodyBytes))
				requestParamsBytes, err := utils.MarshalData(job.RequestParams)
				if err != nil {
					errorChan <- err
					return
				}
				requestQueryBytes, err := utils.MarshalData(job.RequestQuery)
				if err != nil {
					errorChan <- err
					return
				}
				requestBodyBytes, err := utils.MarshalData(job.RequestBody)
				if err != nil {
					errorChan <- err
					return
				}
				requestRoutesChan <- *DTO.NewRouteRequest(job.RequestAuthorization, string(requestQueryBytes), string(requestParamsBytes), string(requestBodyBytes), job.ParentId)
				nextRouteBodyBytes, err := utils.MarshalData(job.NextRouteBody)
				if err != nil {
					errorChan <- err
					return
				}
				nextRouteQueryBytes, err := utils.MarshalData(job.NextRouteQuery)
				if err != nil {
					errorChan <- err
					return
				}
				nextRouteParamsBytes, err := utils.MarshalData(job.NextRouteParams)
				if err != nil {
					errorChan <- err
					return
				}
				nextRoutesChan <- *DTO.NewNextRouteData(string(nextRouteBodyBytes), string(nextRouteQueryBytes),
					string(nextRouteParamsBytes), job.NextAuthorizationHeader, job.ParentId)
				routesInfoChan <- *DTO.NewRouteInfo(job.Path, job.Method, job.ParentId)
			}
		}()
	}
	for _, route := range *routes {
		jobs <- route
	}
	close(jobs)
	wg.Wait()
	close(nextRoutesChan)
	close(requestRoutesChan)
	close(responseRoutesChan)
	close(routesInfoChan)
	select {
	case err := <-errorChan:
		return err
	default:
	}
	var nextRoutes []*DTO.NextRoute
	var requestRoutes []*DTO.RouteRequest
	var responseRoutes []*DTO.RouteResponse
	var routesInfo []*DTO.RouteInfo

	for nextRoute := range nextRoutesChan {
		nextRoutes = append(nextRoutes, &nextRoute)
	}
	for requestRoute := range requestRoutesChan {
		requestRoutes = append(requestRoutes, &requestRoute)
	}
	for responseRoute := range responseRoutesChan {
		responseRoutes = append(responseRoutes, &responseRoute)
	}
	for routeInfo := range routesInfoChan {
		routesInfo = append(routesInfo, &routeInfo)
	}
	nextRoutes = utils.InsertionSortForRoutes(nextRoutes)
	requestRoutes = utils.InsertionSortForRoutes(requestRoutes)
	responseRoutes = utils.InsertionSortForRoutes(responseRoutes)
	routesInfo = utils.InsertionSortForRoutes(routesInfo)
	var routesInfoErr, routesRequestsErr, routesResponsesErr, nextRoutesDataErr error
	var routesInfoIds, routesRequestsIds, routesResponsesIds, nextRoutesDataIds []int
	wg.Add(4)
	go func() {
		defer wg.Done()
		routesInfoIds, routesInfoErr = rs.RouteRepository.InsertRoutesInfo(ctx, routesInfo)
	}()
	go func() {
		defer wg.Done()
		nextRoutesDataIds, nextRoutesDataErr = rs.RouteRepository.InsertNextRoutesData(ctx, nextRoutes)
	}()
	go func() {
		defer wg.Done()
		routesResponsesIds, routesResponsesErr = rs.RouteRepository.InsertRoutesResponses(ctx, responseRoutes)
	}()
	go func() {
		defer wg.Done()
		routesRequestsIds, routesRequestsErr = rs.RouteRepository.InsertRoutesRequests(ctx, requestRoutes)
	}()
	wg.Wait()
	var workingRoutes []DTO.WorkingRoute
	for _, val := range *routes {
		workingRoutes = append(workingRoutes, DTO.WorkingRoute{ParentID: val.ParentId, AppId: appId, Name: name})
	}
	if routesInfoErr != nil {
		return routesInfoErr
	}
	if routesRequestsErr != nil {
		return routesRequestsErr
	}
	if routesResponsesErr != nil {
		return routesResponsesErr
	}
	if nextRoutesDataErr != nil {
		return nextRoutesDataErr
	}
	parentId := 0
	for i := 0; i < len(workingRoutes); i++ {
		workingRoutes[i].NextRouteDataId = nextRoutesDataIds[i]
		workingRoutes[i].RequestID = routesRequestsIds[i]
		workingRoutes[i].ResponseID = routesResponsesIds[i]
		workingRoutes[i].RouteID = routesInfoIds[i]
		workingRoutes[i].ParentID = parentId
		workingRoutes[i].Status = "unknown"
		res, err := rs.RouteRepository.InsertWorkingRoute(ctx, workingRoutes[i])
		if err != nil {
			return err
		}
		parentId = res
	}
	return nil
}
