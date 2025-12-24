package servicesApp

import (
	"context"
	"slices"
	"strings"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
)

type RouteStatusService struct {
	routeRepository interfaces.RouteRepository
	loggerService   utils.LoggerService
}

func NewRouteStatusService(routeRepository interfaces.RouteRepository,
	loggerService utils.LoggerService,
) *RouteStatusService {
	return &RouteStatusService{
		routeRepository: routeRepository,
		loggerService:   loggerService,
	}
}

func (rs *RouteStatusService) sortRoutesToTest(routesToTest []models.RouteToTest) map[string][]models.RouteToTest {
	sortedRoutesToTests := make(map[string][]models.RouteToTest, len(routesToTest))
	for _, routeToTest := range routesToTest {
		key := routeToTest.Name + routeToTest.AppID
		if routeToTest.ParentID == 0 {
			sortedRoutesToTests[key] = append([]models.RouteToTest{routeToTest},
				sortedRoutesToTests[key]...)
		} else {
			sortedRoutesToTests[key] = append(sortedRoutesToTests[key], routeToTest)
		}
	}

	return sortedRoutesToTests
}

func (rs *RouteStatusService) addParamsToThePath(path string, params models.JSONMapStringString) string {
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

func (rs *RouteStatusService) prepareRouteDataForTestRequest(route models.RouteToTest) (authorizationHeader string, url string, preparedBody []byte, err error) {
	authorizationHeader = "Bearer " + route.RequestAuthorization
	query := make([]string, 0, len(route.RequestQuery))
	for key, val := range route.RequestQuery {
		query = append(query, key+"="+val)
	}

	path := rs.addParamsToThePath(route.Path, route.RequestParams)
	url = "http://" + route.IPAddress + ":" + route.Port + path + "?" + strings.Join(query, "&")

	preparedBody, err = utils.MarshalData(route.RequestBody)
	if err != nil {
		rs.loggerService.Error("failed to marshal webhook payload", err)
		return "", "", nil, err
	}

	return authorizationHeader, url, preparedBody, nil
}

func (rs *RouteStatusService) prepareDataForTheNextRoute(route models.RouteToTest, key string,
	val any,
) (nextRouteBody map[string]any, nextRouteParams map[string]string, nextRouteQuery map[string]string, nextRouteAuthorizationHeader string, routeStatus string) {
	routeStatus = "unknown"
	nextRouteBody = make(map[string]any, len(route.RequestBody))
	nextRouteParams = make(map[string]string, len(route.RequestBody))
	nextRouteQuery = make(map[string]string, len(route.RequestBody))
	nextRouteAuthorizationHeader = ""

	if slices.Contains(route.NextRouteBody, key) {
		nextRouteBody[key] = val
	}

	if slices.Contains(route.NextRouteParams, key) {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for param"
			return nil, nil, nil, "", routeStatus
		}
		nextRouteParams[key] = valueConvertedToString
	}

	if slices.Contains(route.NextRouteQuery, key) {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for query"
			return nil, nil, nil, "", routeStatus
		}
		nextRouteQuery[key] = valueConvertedToString
	}
	if route.NextAuthorizationHeader == key {
		valueConvertedToString, ok := val.(string)
		if !ok {
			routeStatus = "Failed;Wrong type of the property for authorization header"
			return nil, nil, nil, "", routeStatus
		}
		nextRouteAuthorizationHeader = valueConvertedToString
	}

	return nextRouteBody, nextRouteParams, nextRouteQuery, nextRouteAuthorizationHeader, routeStatus
}

func (rs *RouteStatusService) CheckRoutesStatus(ctx context.Context) error {
	rs.loggerService.Info("started checking statuses of the routes")

	routesToTest, err := rs.routeRepository.GetWorkingRoutesToTest(ctx)
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

			responseStatusCode, responseBody, err := request.SendHTTP(ctx, url, authorizationHeader, route.Method,
				body, true)
			if err != nil {
				rs.loggerService.Info("Failed to check route", map[string]any{
					"url":    url,
					"method": route.Method,
					"body":   body,
				})
				routeStatus = "Failed;To check route"
				routesStatuses[route.ID] = routeStatus
				break
			}

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
	rs.loggerService.Info("the routes statuses have started inserting into database", routesStatuses)

	err = rs.routeRepository.UpdateWorkingRoutesStatuses(ctx, routesStatuses)
	if err != nil {
		return err
	}

	rs.loggerService.Info("the route statuses have finished inserting into the database.", routesStatuses)

	return nil
}
