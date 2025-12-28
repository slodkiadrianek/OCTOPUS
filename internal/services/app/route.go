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

func (rs *RouteService) CheckRouteStatus(ctx context.Context, routeID int) (string, error) {
	routeStatus, err := rs.routeRepository.CheckRouteStatus(ctx, routeID)
	if err != nil {
		return "", err
	}
	return routeStatus, nil
}

func (rs *RouteService) prepareDataAboutRouteToInsertToDb(routes *[]DTO.CreateRoute) ([]*DTO.NextRoute, []*DTO.RouteRequest, []*DTO.RouteResponse, []*DTO.RouteInfo, error) {
	nextRoutesChan := make(chan DTO.NextRoute, len(*routes))
	requestRoutesChan := make(chan DTO.RouteRequest, len(*routes))
	responseRoutesChan := make(chan DTO.RouteResponse, len(*routes))
	routesInfoChan := make(chan DTO.RouteInfo, len(*routes))
	errorChan := make(chan error, len(*routes))

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
					continue
				}

				responseRoutesChan <- *DTO.NewRouteResponse(job.ResponseStatusCode, job.ParentID, string(responseBodyBytes))

				requestParamsBytes, err := utils.MarshalData(job.RequestParams)
				if err != nil {
					errorChan <- err
					continue
				}

				requestQueryBytes, err := utils.MarshalData(job.RequestQuery)
				if err != nil {
					errorChan <- err
					continue
				}

				requestBodyBytes, err := utils.MarshalData(job.RequestBody)
				if err != nil {
					errorChan <- err
					continue
				}
				requestRoutesChan <- *DTO.NewRouteRequest(job.RequestAuthorization, string(requestQueryBytes), string(requestParamsBytes), string(requestBodyBytes), job.ParentID)

				nextRouteBodyBytes, err := utils.MarshalData(job.NextRouteBody)
				if err != nil {
					errorChan <- err
					continue
				}

				nextRouteQueryBytes, err := utils.MarshalData(job.NextRouteQuery)
				if err != nil {
					errorChan <- err
					continue
				}

				nextRouteParamsBytes, err := utils.MarshalData(job.NextRouteParams)
				if err != nil {
					errorChan <- err
					continue
				}

				nextRoutesChan <- *DTO.NewNextRouteData(string(nextRouteBodyBytes), string(nextRouteQueryBytes),
					string(nextRouteParamsBytes), job.NextRouteAuthorizationHeader, job.ParentID)
				routesInfoChan <- *DTO.NewRouteInfo(job.Path, job.Method, job.ParentID)
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
	close(errorChan)

	select {
	case err := <-errorChan:
		if err != nil {
			return nil, nil, nil, nil, err
		}
	default:
	}

	nextRoutes := make([]*DTO.NextRoute, 0, len(*routes))
	requestRoutes := make([]*DTO.RouteRequest, 0, len(*routes))
	responseRoutes := make([]*DTO.RouteResponse, 0, len(*routes))
	routesInfo := make([]*DTO.RouteInfo, 0, len(*routes))

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
	requestRoutes []*DTO.RouteRequest, responseRoutes []*DTO.RouteResponse, routesInfo []*DTO.RouteInfo,
) (routesInfoIDs, routesRequestsIDs, routesResponsesIDs, nextRoutesDataIDs []int, err error) {
	var wg sync.WaitGroup

	nextRoutes = utils.InsertionSortForRoutes(nextRoutes)
	requestRoutes = utils.InsertionSortForRoutes(requestRoutes)
	responseRoutes = utils.InsertionSortForRoutes(responseRoutes)
	routesInfo = utils.InsertionSortForRoutes(routesInfo)

	var routesInfoErr, routesRequestsErr, routesResponsesErr, nextRoutesDataErr error

	wg.Add(4)
	go func() {
		defer wg.Done()
		routesInfoIDs, routesInfoErr = rs.routeRepository.InsertRoutesInfo(ctx, routesInfo)
	}()
	go func() {
		defer wg.Done()
		nextRoutesDataIDs, nextRoutesDataErr = rs.routeRepository.InsertNextRoutesData(ctx, nextRoutes)
	}()
	go func() {
		defer wg.Done()
		routesResponsesIDs, routesResponsesErr = rs.routeRepository.InsertRoutesResponses(ctx, responseRoutes)
	}()
	go func() {
		defer wg.Done()
		routesRequestsIDs, routesRequestsErr = rs.routeRepository.InsertRoutesRequests(ctx, requestRoutes)
	}()
	wg.Wait()

	if routesInfoErr != nil {
		return nil, nil, nil, nil, routesInfoErr
	}
	if routesRequestsErr != nil {
		return nil, nil, nil, nil, routesRequestsErr
	}
	if routesResponsesErr != nil {
		return nil, nil, nil, nil, routesResponsesErr
	}
	if nextRoutesDataErr != nil {
		return nil, nil, nil, nil, nextRoutesDataErr
	}

	return routesInfoIDs, routesRequestsIDs, routesResponsesIDs, nextRoutesDataIDs, nil
}

func (rs *RouteService) saveWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appID, name string,
	nextRoutesDataIDs,
	routesRequestsIDs, routesResponsesIDs, routesInfoIDs []int,
) error {
	workingRoutes := make([]DTO.WorkingRoute, len(*routes))
	for i, val := range *routes {
		workingRoutes[i] = DTO.WorkingRoute{ParentID: val.ParentID, AppID: appID, Name: name}
	}
	parentID := 0

	for i := 0; i < len(workingRoutes); i++ {
		workingRoutes[i].NextRouteDataID = nextRoutesDataIDs[i]
		workingRoutes[i].RequestID = routesRequestsIDs[i]
		workingRoutes[i].ResponseID = routesResponsesIDs[i]
		workingRoutes[i].RouteID = routesInfoIDs[i]
		workingRoutes[i].ParentID = parentID
		workingRoutes[i].Status = "unknown"
		res, err := rs.routeRepository.InsertWorkingRoute(ctx, workingRoutes[i])
		if err != nil {
			return err
		}
		parentID = res
	}

	return nil
}

func (rs *RouteService) AddWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appID,
	name string,
) error {
	nextRoutes, requestRoutes, responseRoutes, routesInfo, err := rs.prepareDataAboutRouteToInsertToDb(routes)
	if err != nil {
		return err
	}

	nextRoutes = utils.InsertionSortForRoutes(nextRoutes)
	requestRoutes = utils.InsertionSortForRoutes(requestRoutes)
	responseRoutes = utils.InsertionSortForRoutes(responseRoutes)
	routesInfo = utils.InsertionSortForRoutes(routesInfo)

	routesInfoIDs, routesRequestsIDs, routesResponsesIDs, nextRoutesDataIDs, err := rs.saveRouteComponents(ctx, nextRoutes,
		requestRoutes, responseRoutes, routesInfo)
	if err != nil {
		return err
	}

	err = rs.saveWorkingRoutes(ctx, routes, appID, name, nextRoutesDataIDs, routesRequestsIDs,
		routesResponsesIDs, routesInfoIDs)
	if err != nil {
		return err
	}

	return nil
}
