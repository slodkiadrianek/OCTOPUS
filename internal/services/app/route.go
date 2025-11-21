package servicesApp

import (
	"context"
	"runtime"
	"sync"

	"github.com/slodkiadrianek/octopus/internal/services/interfaces"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RouteService struct {
	logger          utils.LoggerService
	routeRepository interfaces.RouteRepository
}

func NewRouteService(logger utils.LoggerService, routeRepository interfaces.RouteRepository) *RouteService {
	return &RouteService{
		logger:          logger,
		routeRepository: routeRepository,
	}
}

func (rs *RouteService) prepareDataAboutRouteToInsertToDb(routes *[]DTO.CreateRoute) ([]*DTO.NextRoute, []*DTO.RouteRequest, []*DTO.RouteResponse, []*DTO.RouteInfo, error) {
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
	wg.Wait()
	close(jobs)
	close(nextRoutesChan)
	close(requestRoutesChan)
	close(responseRoutesChan)
	close(routesInfoChan)
	close(errorChan)
	select {
	case err := <-errorChan:
		return []*DTO.NextRoute{}, []*DTO.RouteRequest{}, []*DTO.RouteResponse{}, []*DTO.RouteInfo{}, err
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
	return nextRoutes, requestRoutes, responseRoutes, routesInfo, nil
}

func (rs *RouteService) saveRouteComponents(ctx context.Context, nextRoutes []*DTO.NextRoute,
	requestRoutes []*DTO.RouteRequest, responseRoutes []*DTO.RouteResponse, routesInfo []*DTO.RouteInfo) ([]int, []int, []int, []int, error) {
	var wg sync.WaitGroup
	nextRoutes = utils.InsertionSortForRoutes(nextRoutes)
	requestRoutes = utils.InsertionSortForRoutes(requestRoutes)
	responseRoutes = utils.InsertionSortForRoutes(responseRoutes)
	routesInfo = utils.InsertionSortForRoutes(routesInfo)
	var routesInfoErr, routesRequestsErr, routesResponsesErr, nextRoutesDataErr error
	var routesInfoIds, routesRequestsIds, routesResponsesIds, nextRoutesDataIds []int
	wg.Add(4)
	go func() {
		defer wg.Done()
		routesInfoIds, routesInfoErr = rs.routeRepository.InsertRoutesInfo(ctx, routesInfo)
	}()
	go func() {
		defer wg.Done()
		nextRoutesDataIds, nextRoutesDataErr = rs.routeRepository.InsertNextRoutesData(ctx, nextRoutes)
	}()
	go func() {
		defer wg.Done()
		routesResponsesIds, routesResponsesErr = rs.routeRepository.InsertRoutesResponses(ctx, responseRoutes)
	}()
	go func() {
		defer wg.Done()
		routesRequestsIds, routesRequestsErr = rs.routeRepository.InsertRoutesRequests(ctx, requestRoutes)
	}()
	wg.Wait()

	if routesInfoErr != nil {
		return []int{}, []int{}, []int{}, []int{}, routesInfoErr
	}
	if routesRequestsErr != nil {
		return []int{}, []int{}, []int{}, []int{}, routesRequestsErr
	}
	if routesResponsesErr != nil {
		return []int{}, []int{}, []int{}, []int{}, routesResponsesErr
	}
	if nextRoutesDataErr != nil {
		return []int{}, []int{}, []int{}, []int{}, nextRoutesDataErr
	}
	return routesInfoIds, routesRequestsIds, routesResponsesIds, nextRoutesDataIds, nil
}
func (rs *RouteService) saveWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appId, name string,
	nextRoutesDataIds,
	routesRequestsIds, routesResponsesIds, routesInfoIds []int) error {
	var workingRoutes []DTO.WorkingRoute
	for _, val := range *routes {
		workingRoutes = append(workingRoutes, DTO.WorkingRoute{ParentID: val.ParentId, AppId: appId, Name: name})
	}
	parentId := 0
	for i := 0; i < len(workingRoutes); i++ {
		workingRoutes[i].NextRouteDataId = nextRoutesDataIds[i]
		workingRoutes[i].RequestID = routesRequestsIds[i]
		workingRoutes[i].ResponseID = routesResponsesIds[i]
		workingRoutes[i].RouteID = routesInfoIds[i]
		workingRoutes[i].ParentID = parentId
		workingRoutes[i].Status = "unknown"
		res, err := rs.routeRepository.InsertWorkingRoute(ctx, workingRoutes[i])
		if err != nil {
			return err
		}
		parentId = res
	}
	return nil
}
func (rs *RouteService) AddWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appId,
	name string) error {
	nextRoutes, requestRoutes, responseRoutes, routesInfo, err := rs.prepareDataAboutRouteToInsertToDb(routes)
	if err != nil {
		return err
	}
	nextRoutes = utils.InsertionSortForRoutes(nextRoutes)
	requestRoutes = utils.InsertionSortForRoutes(requestRoutes)
	responseRoutes = utils.InsertionSortForRoutes(responseRoutes)
	routesInfo = utils.InsertionSortForRoutes(routesInfo)
	routesInfoIds, routesRequestsIds, routesResponsesIds, nextRoutesDataIds, err := rs.saveRouteComponents(ctx, nextRoutes,
		requestRoutes, responseRoutes, routesInfo)
	if err != nil {
		return err
	}
	err = rs.saveWorkingRoutes(ctx, routes, appId, name, nextRoutesDataIds, routesRequestsIds,
		routesResponsesIds, routesInfoIds)
	return nil
}
