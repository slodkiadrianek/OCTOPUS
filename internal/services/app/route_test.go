package servicesApp

import (
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRouteService_prepareDataAboutRouteToInsertToDb(t *testing.T) {
	type args struct {
		name                  string
		routes                *[]DTO.CreateRoute
		expectedNextRoutes    []*DTO.NextRoute
		expectedRouteRequest  []*DTO.RouteRequest
		expectedRouteResponse []*DTO.RouteResponse
		expectedRouteInfo     []*DTO.RouteInfo
		expectedError         error
	}

	testsScenarios := []args{
		{
			name: "Properly prepared routes to insert to db",
			routes: &[]DTO.CreateRoute{
				{
					ResponseBody: map[string]any{
						"id":     1,
						"body":   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
						"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
						"userId": 1,
					},
					ParentId: 0,
					RequestParams: map[string]string{
						"postId": "1",
					},
					RequestQuery:            map[string]string{},
					RequestBody:             map[string]any{},
					RequestAuthorization:    "",
					NextRouteBody:           make([]string, 0),
					NextRouteParams:         make([]string, 0),
					NextRouteQuery:          make([]string, 0),
					NextAuthorizationHeader: "",
					Path:                    "/posts/{postId}",
					Method:                  "GET",
				},
			},
			expectedNextRoutes: []*DTO.NextRoute{
				{
					Body:                "",
					Params:              "",
					Query:               "",
					AuthorizationHeader: "",
				},
			},
			expectedRouteRequest: []*DTO.RouteRequest{
				{
					AuthorizationHeader: "",
					Query:               "",
					Params:              "",
				},
			},
			expectedRouteResponse: []*DTO.RouteResponse{
				{
					Body: "",
				},
			}
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			routeRepository := new(mocks.MockRouteRepository)
			routeService := NewRouteService(loggerService, routeRepository)
			nextRoutes, requestRoutes, responseRoutes, routesInfo, err := routeService.prepareDataAboutRouteToInsertToDb(testScenario.routes)
			assert.Equal(t, testScenario.expectedNextRoutes, nextRoutes)
			assert.Equal(t, testScenario.expectedRouteRequest, requestRoutes)
			assert.Equal(t, testScenario.expectedRouteResponse, responseRoutes)
			as
			assert.Equal(t, testScenario.expectedError, err)
		})
	}
}
