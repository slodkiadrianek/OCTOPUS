package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RouteService struct {
	Logger          *utils.Logger
	RouteRepository *repository.RouteRepository
}

func NewRouteService(logger *utils.Logger, routeRepository *repository.RouteRepository) *RouteService {
	return &RouteService{
		Logger:          logger,
		RouteRepository: routeRepository,
	}
}

func (rs *RouteService) AddWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appId string, name string) error {
	nextRoutesDataChan := make(chan DTO.NextRouteData, len(*routes))
	requestRoutesDataChan := make(chan DTO.RouteRequest, len(*routes))
	responseRoutesDataChan := make(chan DTO.RouteResponse, len(*routes))
	routesDataChan := make(chan DTO.RouteInfo, len(*routes))
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
				responseRoutesDataChan <- *DTO.NewRouteResponse(job.ResponseStatusCode, job.ParentId, string(responseBodyBytes))
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
				requestRoutesDataChan <- *DTO.NewRouteRequest(job.RequestAuthorization, string(requestQueryBytes), string(requestParamsBytes), string(requestBodyBytes), job.ParentId)
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
				nextRoutesDataChan <- *DTO.NewNextRouteData(string(nextRouteBodyBytes), string(nextRouteQueryBytes), string(nextRouteParamsBytes), job.NextAuthorizationHeader)
				routesDataChan <- *DTO.NewRouteInfo(job.Path, job.Method, job.ParentId)
			}
		}()
	}
	for _, route := range *routes {
		jobs <- route
	}
	close(jobs)
	wg.Wait()
	close(nextRoutesDataChan)
	close(requestRoutesDataChan)
	close(responseRoutesDataChan)
	close(routesDataChan)
	select {
	case err := <-errorChan:
		return err
	default:
	}
	var nextRoutesData []*DTO.NextRouteData
	var requestRoutesData []*DTO.RouteRequest
	var responseRoutesData []*DTO.RouteResponse
	var routesInfoData []*DTO.RouteInfo

	for data := range nextRoutesDataChan {
		nextRoutesData = append(nextRoutesData, &data)
	}
	for data := range requestRoutesDataChan {
		requestRoutesData = append(requestRoutesData, &data)
	}
	for data := range responseRoutesDataChan {
		responseRoutesData = append(responseRoutesData, &data)
	}
	for data := range routesDataChan {
		routesInfoData = append(routesInfoData, &data)
	}
	nextRoutesData = utils.InsertionSortForRoutes(nextRoutesData)
	requestRoutesData = utils.InsertionSortForRoutes(requestRoutesData)
	responseRoutesData = utils.InsertionSortForRoutes(responseRoutesData)
	routesInfoData = utils.InsertionSortForRoutes(routesInfoData)
	var routesInfoErr, routesRequestsErr, routesResponsesErr, nextRoutesDataErr error
	var routesInfoIds, routesRequestsIds, routesResponsesIds, nextRoutesDataIds []int
	wg.Add(4)
	go func() {
		defer wg.Done()
		routesInfoIds, routesInfoErr = rs.RouteRepository.InsertRoutesInfo(ctx, routesInfoData)
	}()
	go func() {
		defer wg.Done()
		nextRoutesDataIds, nextRoutesDataErr = rs.RouteRepository.InsertNextRoutesData(ctx, nextRoutesData)
	}()
	go func() {
		defer wg.Done()
		routesResponsesIds, routesResponsesErr = rs.RouteRepository.InsertRoutesResponses(ctx, responseRoutesData)
	}()
	go func() {
		defer wg.Done()
		routesRequestsIds, routesRequestsErr = rs.RouteRepository.InsertRoutesReuquests(ctx, requestRoutesData)
	}()
	wg.Wait()
	var workingRoutes []DTO.WorkingRoute
	for _, val := range *routes {
		workingRoutes = append(workingRoutes, DTO.WorkingRoute{ParentId: val.ParentId, AppId: appId, Name: name})
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
	fmt.Println(nextRoutesDataIds)
	for i := 0; i < len(workingRoutes); i++ {
		workingRoutes[i].NextRouteDataId = nextRoutesDataIds[i]
		workingRoutes[i].RequestId = routesRequestsIds[i]
		workingRoutes[i].ResponseId = routesResponsesIds[i]
		workingRoutes[i].RouteId = routesInfoIds[i]
		workingRoutes[i].ParentId = parentId
		workingRoutes[i].Status = "unknown"
		res, err := rs.RouteRepository.InsertWorkingRoute(ctx, workingRoutes[i])
		if err != nil {
			return err
		}
		parentId = res

	}
	return nil
}
