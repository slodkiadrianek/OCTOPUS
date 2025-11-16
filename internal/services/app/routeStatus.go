package servicesApp

import (
	"context"
	"slices"
	"strings"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RouteStatusService struct {
	RouteRepository interfaces.RouteRepository
	LoggerService   utils.Logger
}

func NewRouteStatusService(routeRepository interfaces.RouteRepository,
	loggerService utils.Logger) *RouteStatusService {
	return &RouteStatusService{
		RouteRepository: routeRepository,
		LoggerService:   loggerService,
	}
}

func (rs *RouteStatusService) sortRoutesToTest(routesToTest []models.RouteToTest) map[string][]models.RouteToTest {
	sortedRoutesToTests := make(map[string][]models.RouteToTest)
	for _, routeToTest := range routesToTest {
		key := routeToTest.Name + routeToTest.AppId
		if routeToTest.ParentID == 0 {
			sortedRoutesToTests[key] = append([]models.RouteToTest{routeToTest},
				sortedRoutesToTests[key]...)
		} else {
			sortedRoutesToTests[key] = append(sortedRoutesToTests[key], routeToTest)
		}
	}

	return sortedRoutesToTests
}

func (rs *RouteStatusService) addParamsToThePath(path string, params models.JsonMapStringString) string {
	splittedPath := strings.Split(path, "/")
	for i := 0; i < len(splittedPath); i++ {
		partOfPath := splittedPath[i]
		leftBrace := strings.Contains(partOfPath, "{")
		rightBrace := strings.Contains(partOfPath, "}")

		if leftBrace && rightBrace {
			param := partOfPath[1 : len(partOfPath)-1]
			splittedPath[i] = params[param]
		}
	}
	pathWithParamsIncluded := strings.Join(splittedPath, "/")
	return pathWithParamsIncluded
}

func (rs *RouteStatusService) prepareRouteDataForTestRequest(route models.RouteToTest) (string, string, []byte, error) {
	authorizationHeader := "Bearer " + route.RequestAuthorization
	var query []string
	for key, val := range route.RequestQuery {
		query = append(query, key+"="+val)
	}

	path := rs.addParamsToThePath(route.Path, route.RequestParams)
	url := "http://" + route.IpAddress + ":" + route.Port + path + "?" + strings.Join(query, "&")

	jsonData, err := utils.MarshalData(route.RequestBody)
	if err != nil {
		rs.LoggerService.Error("Failed to marshal webhook payload", err)
		return "", "", []byte{}, err
	}

	return authorizationHeader, url, jsonData, nil
}

func (rs *RouteStatusService) prepareDataForTheNextRoute(route models.RouteToTest, key string,
	val any,
) (map[string]any, map[string]string, map[string]string, string, string) {
	routeStatus := "unknown"
	nextRouteBody := make(map[string]any)
	nextRouteParams := make(map[string]string)
	nextRouteQuery := make(map[string]string)
	nextRouteAuthorizationHeader := ""

	if slices.Contains(route.NextRouteBody, key) {
		nextRouteBody[key] = val
	}

	if slices.Contains(route.NextRouteParams, key) {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for param"
			return map[string]any{}, map[string]string{}, map[string]string{}, "", routeStatus
		}
		nextRouteParams[key] = valueConvertedToString
	}

	if slices.Contains(route.NextRouteQuery, key) {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for query"
			return map[string]any{}, map[string]string{}, map[string]string{}, "", routeStatus
		}
		nextRouteQuery[key] = valueConvertedToString
	}

	valueConvertedToString, ok := val.(string)
	if !ok {
		routeStatus = "Failed;Wrong type of the property for authorization header"
		return map[string]any{}, map[string]string{}, map[string]string{}, "", routeStatus
	}
	if strings.Contains(valueConvertedToString, "eyJlbWFpbCI6IlRFU1QiLCJleHAiOjE3N") {
		nextRouteAuthorizationHeader = valueConvertedToString
	}

	return nextRouteBody, nextRouteParams, nextRouteQuery, nextRouteAuthorizationHeader, routeStatus
}

func (rs *RouteStatusService) CheckRoutesStatus(ctx context.Context) error {
	rs.LoggerService.Info("Started checking statuses of the routes")

	routesToTest, err := rs.RouteRepository.GetWorkingRoutesToTest(ctx)
	if err != nil {
		return err
	}

	if len(routesToTest) < 1 {
		return nil
	}

	sortedRoutesToTests := rs.sortRoutesToTest(routesToTest)

	routesStatuses := make(map[int]string)
	for _, routesToTest := range sortedRoutesToTests {
		nextRouteBody := make(map[string]any)
		nextRouteParams := make(map[string]string)
		nextRouteQuery := make(map[string]string)
		nextRouteAuthorizationHeader := ""

		for _, route := range routesToTest {
			routeStatus := "unknown"

			if len(nextRouteBody) > 0 {
				route.RequestBody = nextRouteBody
			}

			if len(nextRouteParams) > 0 {
				route.RequestParams = nextRouteParams
			}

			if len(nextRouteQuery) > 0 {
				route.RequestQuery = nextRouteQuery
			}

			if len(nextRouteAuthorizationHeader) > 0 {
				route.RequestAuthorization = nextRouteAuthorizationHeader
			}

			authorizationHeader, url, body, err := rs.prepareRouteDataForTestRequest(route)
			if err != nil {
				return err
			}

			responseStatusCode, responseBody, err := utils.DoHttpRequest(ctx, url, authorizationHeader, route.Method, body, rs.LoggerService)
			if len(responseBody) != len(route.ResponseBody) {
				routeStatus = "Failed;Different body"
				routesStatuses[route.ID] = routeStatus
				break
			}
			if responseStatusCode != route.ResponseStatusCode {
				routeStatus = "Failed;Status Code"
				routesStatuses[route.ID] = routeStatus
				break
			}

			for key, val := range responseBody {
				nextRouteBody, nextRouteParams, nextRouteQuery, nextRouteAuthorizationHeader,
					routeStatus = rs.prepareDataForTheNextRoute(route, key, val)
			}

			routeStatus = "success"
			routesStatuses[route.ID] = routeStatus
		}

	}
	rs.LoggerService.Info("The routes statuses have started inserting into database", routesStatuses)

	err = rs.RouteRepository.UpdateWorkingRoutesStatuses(ctx, routesStatuses)
	if err != nil {
		return err
	}

	rs.LoggerService.Info("The route statuses have finished inserting into the database.", routesStatuses)
	return nil
}
