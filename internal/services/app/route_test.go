package servicesApp

import (
	"context"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
					Body:   "[]",
					Query:  "[]",
					Params: "[]",
				},
			},
			expectedRouteRequest: []*DTO.RouteRequest{
				{
					Params: "{\"postId\":\"1\"}",
					Query:  "{}",
					Body:   "{}",
				},
			},
			expectedRouteResponse: []*DTO.RouteResponse{
				{
					Body: "{\"body\":\"quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto\",\"id\":1,\"title\":\"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\"userId\":1}",
				},
			},
			expectedRouteInfo: []*DTO.RouteInfo{
				{
					Path:   "/posts/{postId}",
					Method: "GET",
				},
			},
			expectedError: nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			routeRepository := new(mocks.MockRouteRepository)
			routeService := NewRouteService(loggerService, routeRepository)
			nextRoutes, requestRoutes, responseRoutes, routeInfo, err := routeService.prepareDataAboutRouteToInsertToDb(testScenario.routes)
			assert.Equal(t, testScenario.expectedNextRoutes, nextRoutes)
			assert.Equal(t, testScenario.expectedRouteRequest, requestRoutes)
			assert.Equal(t, testScenario.expectedRouteResponse, responseRoutes)
			assert.Equal(t, testScenario.expectedRouteInfo, routeInfo)
			assert.Equal(t, testScenario.expectedError, err)
		})
	}
}

func TestRouteService_saveRouteComponents(t *testing.T) {
	type args struct {
		name                       string
		nextRoutes                 []*DTO.NextRoute
		routeRequest               []*DTO.RouteRequest
		routeResponse              []*DTO.RouteResponse
		routeInfo                  []*DTO.RouteInfo
		expectedRoutesInfoIds      []int
		expectedRouteRequestsIds   []int
		expectedRoutesResponsesIds []int
		expectedNextRouteDataIds   []int
		expectedError              error
		setupMocks                 func() interfaces.RouteRepository
	}

	testsScenarios := []args{
		{
			name: "Properly saved components in db",
			nextRoutes: []*DTO.NextRoute{
				{
					Body:   "[]",
					Query:  "[]",
					Params: "[]",
				},
			},
			routeRequest: []*DTO.RouteRequest{
				{
					Params: "{\"postId\":\"1\"}",
					Query:  "{}",
					Body:   "{}",
				},
			},
			routeResponse: []*DTO.RouteResponse{
				{
					Body: "{\"body\":\"quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto\",\"id\":1,\"title\":\"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\"userId\":1}",
				},
			},
			routeInfo: []*DTO.RouteInfo{
				{
					Path:   "/posts/{postId}",
					Method: "GET",
				},
			},
			setupMocks: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("InsertRoutesInfo", mock.Anything, mock.Anything).Return([]int{0}, nil)
				mRouteRepository.On("InsertNextRoutesData", mock.Anything, mock.Anything).Return([]int{0}, nil)

				mRouteRepository.On("InsertRoutesResponses", mock.Anything, mock.Anything).Return([]int{0}, nil)
				mRouteRepository.On("InsertRoutesRequests", mock.Anything, mock.Anything).Return([]int{0}, nil)
				return mRouteRepository
			},
			expectedError:              nil,
			expectedRouteRequestsIds:   []int{0},
			expectedRoutesResponsesIds: []int{0},
			expectedRoutesInfoIds:      []int{0},
			expectedNextRouteDataIds:   []int{0},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			routeRepository := testScenario.setupMocks()
			routeService := NewRouteService(loggerService, routeRepository)
			routesInfoIds, routesRequestsIds, routesResponsesIds, nextRoutesDataIds, err := routeService.saveRouteComponents(ctx, testScenario.nextRoutes, testScenario.routeRequest, testScenario.routeResponse, testScenario.routeInfo)
			assert.Equal(t, testScenario.expectedNextRouteDataIds, nextRoutesDataIds)
			assert.Equal(t, testScenario.expectedError, err)
			assert.Equal(t, testScenario.expectedRoutesInfoIds, routesInfoIds)
			assert.Equal(t, testScenario.expectedRoutesResponsesIds, routesResponsesIds)
			assert.Equal(t, testScenario.expectedRouteRequestsIds, routesRequestsIds)
		})
	}
}

func TestRouteService_saveWorkingRoutes(t *testing.T) {
	type args struct {
		name                  string
		routes                *[]DTO.CreateRoute
		appId                 string
		nameOfTheWorkingRoute string
		routesInfoIds         []int
		routeRequestsIds      []int
		routesResponsesIds    []int
		nextRouteDataIds      []int
		expectedError         error
		setupMocks            func() interfaces.RouteRepository
	}

	testsScenarios := []args{
		{
			name:                  "Properly saved working routes",
			appId:                 "27cf4966c158762ceb9495fbdd044a73325efd3bd2a4f9646fc45662ef59490d",
			nameOfTheWorkingRoute: "test",
			routes: &[]DTO.CreateRoute{{
				ParentId: 0,
			}},
			expectedError:      nil,
			routeRequestsIds:   []int{0},
			routesResponsesIds: []int{0},
			routesInfoIds:      []int{0},
			nextRouteDataIds:   []int{0},
			setupMocks: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("InsertWorkingRoute", mock.Anything, mock.Anything).Return(0, nil)
				return mRouteRepository
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			routeRepository := testScenario.setupMocks()
			routeService := NewRouteService(loggerService, routeRepository)
			err := routeService.saveWorkingRoutes(ctx, testScenario.routes, testScenario.appId, testScenario.nameOfTheWorkingRoute, testScenario.nextRouteDataIds, testScenario.routeRequestsIds, testScenario.routesResponsesIds, testScenario.routesInfoIds)
			assert.Equal(t, testScenario.expectedError, err)
		})
	}
}

func TestRouteService_AddWorkingRoutes(t *testing.T) {
	type args struct {
		name                  string
		routes                *[]DTO.CreateRoute
		appId                 string
		nameOfTheWorkingRoute string
		expectedError         error
		setupMocks            func() interfaces.RouteRepository
	}

	testsScenarios := []args{{
		name: "Properly added working routes to the database",
		routes: &[]DTO.CreateRoute{
			{
				Path:               "/user",
				Method:             "GET",
				ResponseStatusCode: 200,
			},
		},
		appId:                 "27cf4966c158762ceb9495fbdd044a73325efd3bd2a4f9646fc45662ef59490d",
		nameOfTheWorkingRoute: "test",
		expectedError:         nil,
		setupMocks: func() interfaces.RouteRepository {
			mRouteRepository := new(mocks.MockRouteRepository)
			mRouteRepository.On("InsertRoutesInfo", mock.Anything, mock.Anything).Return([]int{0}, nil)
			mRouteRepository.On("InsertNextRoutesData", mock.Anything, mock.Anything).Return([]int{0}, nil)

			mRouteRepository.On("InsertRoutesResponses", mock.Anything, mock.Anything).Return([]int{0}, nil)
			mRouteRepository.On("InsertRoutesRequests", mock.Anything, mock.Anything).Return([]int{0}, nil)
			mRouteRepository.On("InsertWorkingRoute", mock.Anything, mock.Anything).Return(0, nil)
			return mRouteRepository
		},
	}}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			routeRepository := testScenario.setupMocks()
			routeService := NewRouteService(loggerService, routeRepository)
			err := routeService.AddWorkingRoutes(ctx, testScenario.routes, testScenario.appId, testScenario.nameOfTheWorkingRoute)
			assert.Equal(t, testScenario.expectedError, err)
		})
	}

}
