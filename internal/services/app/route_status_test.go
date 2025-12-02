package servicesApp

import (
	"context"
	"errors"

	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
)

func TestSortRoutesToTest(t *testing.T) {
	type args struct {
		name         string
		routeToTest  []models.RouteToTest
		expectedData map[string][]models.RouteToTest
	}
	testsScenarios := []args{
		{
			name: "Properly sorted routes",
			routeToTest: []models.RouteToTest{
				{
					ParentID: 0,
					Name:     "First Route",
					AppId:    "1",
				},
				{
					ParentID: 1,
					Name:     "First Route",
					AppId:    "1",
				},
				{
					ParentID: 0,
					Name:     "Third Route",
					AppId:    "3",
				},
			},
			expectedData: map[string][]models.RouteToTest{
				"First Route1": {
					{
						ParentID: 0,
						Name:     "First Route",
						AppId:    "1",
					},
					{
						ParentID: 1,
						Name:     "First Route",
						AppId:    "1",
					},
				},
				"Third Route3": {
					{
						ParentID: 0,
						Name:     "Third Route",
						AppId:    "3",
					},
				},
			},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			env, err := config.SetConfig("../../../.env")
			if err != nil {
				panic(err)
			}
			db, err := config.NewDb(env.DbLink, "postgres")
			if err != nil {
				panic(err)
			}
			routeRepository := repository.NewRouteRepository(db.DbConnection, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			sortedData := routeStatusService.sortRoutesToTest(testScenario.routeToTest)
			assert.Equal(t, testScenario.expectedData, sortedData)
		})
	}
}

func TestAddParamsToThePath(t *testing.T) {
	type args struct {
		name         string
		path         string
		params       models.JsonMapStringString
		expectedData string
	}
	testsScenarios := []args{
		{
			name: "Properly added params to the path",
			path: "/users/{userId}",
			params: models.JsonMapStringString{
				"userId": "userId",
			},
			expectedData: "/users/userId",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			env, err := config.SetConfig("../../../.env")
			if err != nil {
				panic(err)
			}
			db, err := config.NewDb(env.DbLink, "postgres")
			if err != nil {
				panic(err)
			}
			routeRepository := repository.NewRouteRepository(db.DbConnection, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			pathWithParamsIncluded := routeStatusService.addParamsToThePath(testScenario.path, testScenario.params)
			assert.Equal(t, testScenario.expectedData, pathWithParamsIncluded)
		})
	}
}

func TestRouteStatusService_prepareDataForTestRequest(t *testing.T) {
	type args struct {
		name                        string
		route                       models.RouteToTest
		expectedAuthorizationHeader string
		expectedUrl                 string
		expectedJsonData            []byte
		expectedError               error
	}
	testsScenarios := []args{
		{
			name: "Properly prepared route for the test request",
			route: models.RouteToTest{
				RequestAuthorization: "fj349f83hf893h9834fh834",
				RequestQuery:         map[string]string{"userId": "userId"},
				Path:                 "/users/{userId}",
				RequestParams: map[string]string{
					"userId": "userId",
				},
				RequestBody: map[string]any{"appId": "appId"},
				IpAddress:   "127.0.0.1",
				Port:        "8080",
			},
			expectedAuthorizationHeader: "Bearer fj349f83hf893h9834fh834",
			expectedUrl:                 "http://127.0.0.1:8080/users/userId?userId=userId",
			expectedJsonData:            []uint8([]byte{0x7b, 0x22, 0x61, 0x70, 0x70, 0x49, 0x64, 0x22, 0x3a, 0x22, 0x61, 0x70, 0x70, 0x49, 0x64, 0x22, 0x7d}),
			expectedError:               nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			env, err := config.SetConfig("../../../.env")
			if err != nil {
				panic(err)
			}
			db, err := config.NewDb(env.DbLink, "postgres")
			if err != nil {
				panic(err)
			}
			routeRepository := repository.NewRouteRepository(db.DbConnection, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			authorizationHeader, url, jsonData, err := routeStatusService.prepareRouteDataForTestRequest(testScenario.route)
			assert.Equal(t, testScenario.expectedAuthorizationHeader, authorizationHeader)
			assert.Equal(t, testScenario.expectedUrl, url)
			assert.Equal(t, testScenario.expectedJsonData, jsonData)
			assert.Equal(t, testScenario.expectedError, err)
		})
	}
}

func TestRouteStatusService_prepareRouteDataForTheNextRoute(t *testing.T) {
	type args struct {
		name                             string
		route                            models.RouteToTest
		key                              string
		val                              any
		expectedNextRouteBody            map[string]any
		expectedNextRouteParams          map[string]string
		expectedRouteQuery               map[string]string
		expectedRouteAuthorizationHeader string
		expectedRouteStatus              string
	}
	testsScenarios := []args{
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteBody",
			route: models.RouteToTest{
				NextRouteBody: models.JsonStringSlice{"userId"},
			},
			key:                              "userId",
			val:                              "userId",
			expectedNextRouteBody:            map[string]any{"userId": "userId"},
			expectedNextRouteParams:          map[string]string{},
			expectedRouteQuery:               map[string]string{},
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteParams",
			route: models.RouteToTest{
				NextRouteParams: models.JsonStringSlice{"userId"},
			},
			key:                              "userId",
			val:                              "userId",
			expectedNextRouteBody:            map[string]any{},
			expectedNextRouteParams:          map[string]string{"userId": "userId"},
			expectedRouteQuery:               map[string]string{},
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteQuery",
			route: models.RouteToTest{
				NextRouteQuery: models.JsonStringSlice{"userId"},
			},
			key:                              "userId",
			val:                              "userId",
			expectedNextRouteBody:            map[string]any{},
			expectedNextRouteParams:          map[string]string{},
			expectedRouteQuery:               map[string]string{"userId": "userId"},
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteAuthorizationHeader",
			route: models.RouteToTest{
				NextAuthorizationHeader: "userId",
			},
			key:                              "userId",
			val:                              "userId",
			expectedNextRouteBody:            map[string]any{},
			expectedNextRouteParams:          map[string]string{},
			expectedRouteQuery:               map[string]string{},
			expectedRouteAuthorizationHeader: "userId",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Failed to assign value to the nexRouteQuery",
			route: models.RouteToTest{
				NextRouteQuery: models.JsonStringSlice{"userId"},
			},
			key:                              "userId",
			val:                              2,
			expectedNextRouteBody:            nil,
			expectedNextRouteParams:          nil,
			expectedRouteQuery:               nil,
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "Failed;Wrong type of the property for query",
		},
		{
			name: "Failed to assign value to the nexRouteParams",
			route: models.RouteToTest{
				NextRouteParams: models.JsonStringSlice{"userId"},
			},
			key:                              "userId",
			val:                              2,
			expectedNextRouteBody:            nil,
			expectedNextRouteParams:          nil,
			expectedRouteQuery:               nil,
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "Failed;Wrong type of the property for param",
		},
		{
			name: "Failed to assign value to the nexRouteAuthorizationHeader",
			route: models.RouteToTest{
				NextAuthorizationHeader: "userId",
			},
			key:                              "userId",
			val:                              2,
			expectedNextRouteBody:            nil,
			expectedNextRouteParams:          nil,
			expectedRouteQuery:               nil,
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "Failed;Wrong type of the property for authorization header",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {

			loggerService := tests.CreateLogger()
			env, err := config.SetConfig("../../../.env")
			if err != nil {
				panic(err)
			}
			db, err := config.NewDb(env.DbLink, "postgres")
			if err != nil {
				panic(err)
			}
			routeRepository := repository.NewRouteRepository(db.DbConnection, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			nextRouteBody, nextRouteParams, nextRouteQuery, nextRouteAuthorizationHeader, routeStatus := routeStatusService.prepareDataForTheNextRoute(testScenario.route, testScenario.key, testScenario.val)
			assert.Equal(t, testScenario.expectedNextRouteBody, nextRouteBody)
			assert.Equal(t, testScenario.expectedNextRouteParams, nextRouteParams)
			assert.Equal(t, testScenario.expectedRouteQuery, nextRouteQuery)
			assert.Equal(t, testScenario.expectedRouteAuthorizationHeader, nextRouteAuthorizationHeader)
			assert.Equal(t, testScenario.expectedRouteStatus, routeStatus)
		})
	}
}

func TestRouteStatusService_CheckRoutesStatus(t *testing.T) {
	type args struct {
		name          string
		setupMock     func() interfaces.RouteRepository
		expectedError error
	}
	testsScenarios := []args{
		{
			name: "Failed to get data from db",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{}, errors.New("Failed to get data from db"))
				return mRouteRepository
			},
			expectedError: errors.New("Failed to get data from db"),
		},
		{
			name: "Lack of routes to test",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{}, nil)
				return mRouteRepository
			},
			expectedError: nil,
		},
		{
			name: "Proper single route test",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{
					{
						IpAddress: "jsonplaceholder.typicode.com",
						Port:      "80",
						Path:      "/posts/{postId}",
						Method:    "GET",
						RequestParams: map[string]string{
							"postId": "1",
						},
						ResponseBody: map[string]any{
							"id":     1,
							"body":   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
							"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
							"userId": 1,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(nil)
				return mRouteRepository
			},
			expectedError: nil,
		},
		{
			name: "Failed to update working route status",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{
					{
						IpAddress: "jsonplaceholder.typicode.com",
						Port:      "80",
						Path:      "/posts/{postId}",
						Method:    "GET",
						RequestParams: map[string]string{
							"postId": "1",
						},
						ResponseBody: map[string]any{
							"id":     1,
							"body":   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
							"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
							"userId": 1,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(errors.New("Failed to update working route status"))
				return mRouteRepository
			},
			expectedError: errors.New("Failed to update working route status"),
		},
		{
			name: "Properly chained 2 routes",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{
					{
						IpAddress: "jsonplaceholder.typicode.com",
						Port:      "80",
						Path:      "/posts/{postId}",
						Method:    "GET",
						RequestParams: map[string]string{
							"postId": "1",
						},
						ResponseBody: map[string]any{
							"id":     1,
							"body":   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
							"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
							"userId": 1,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(errors.New("Failed to update working route status"))
				return mRouteRepository
			},
			expectedError: errors.New("Failed to update working route status"),
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			ctx := context.Background()
			routeRepository := testScenario.setupMock()
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			err := routeStatusService.CheckRoutesStatus(ctx)
			assert.Equal(t, testScenario.expectedError, err)
		})
	}
}
