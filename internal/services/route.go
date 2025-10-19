package services

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RouteService struct {
	Logger *utils.Logger
}

func NewRouteService(logger *utils.Logger) *RouteService {
	return &RouteService{
		Logger: logger,
	}
}

func (rs *RouteService) AddWorkingRoutes(routes *[]DTO.CreateRoute, appId int) error {
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
				nextRoutesDataChan <- *DTO.NewNextRouteData(string(nextRouteBodyBytes), string(nextRouteQueryBytes), string(nextRouteParamsBytes), job.ParentId)
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
	nextRoutesData = utils.InsertionSortForRoutes[*DTO.NextRouteData](nextRoutesData)
	requestRoutesData = utils.InsertionSortForRoutes[*DTO.RouteRequest](requestRoutesData)
	responseRoutesData = utils.InsertionSortForRoutes[*DTO.RouteResponse](responseRoutesData)
	routesInfoData = utils.InsertionSortForRoutes[*DTO.RouteInfo](routesInfoData)

	for _, val := range nextRoutesData {
		fmt.Println(val)
	}
	for _, val := range requestRoutesData {
		fmt.Println(val)
	}
	for _, val := range responseRoutesData {
		fmt.Println(val)
	}
	for _, val := range routesInfoData {
		fmt.Println(val)
	}
	return nil
}
