package servicesApp

import (
	"context"
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteStatusService_sortRoutesToTest(t *testing.T) {
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
					AppID:    "1",
				},
				{
					ParentID: 1,
					Name:     "First Route",
					AppID:    "1",
				},
				{
					ParentID: 0,
					Name:     "Third Route",
					AppID:    "3",
				},
			},
			expectedData: map[string][]models.RouteToTest{
				"First Route1": {
					{
						ParentID: 0,
						Name:     "First Route",
						AppID:    "1",
					},
					{
						ParentID: 1,
						Name:     "First Route",
						AppID:    "1",
					},
				},
				"Third Route3": {
					{
						ParentID: 0,
						Name:     "Third Route",
						AppID:    "3",
					},
				},
			},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			routeRepositoryMock := new(mocks.MockRouteRepository)
			routeStatusService := NewRouteStatusService(routeRepositoryMock, loggerService)
			sortedData := routeStatusService.sortRoutesToTest(testScenario.routeToTest)
			assert.Equal(t, testScenario.expectedData, sortedData)
		})
	}
}

func TestRouteStatusService_addParamsToThePath(t *testing.T) {
	type args struct {
		name         string
		path         string
		params       models.JSONMapStringString
		expectedData string
	}
	testsScenarios := []args{
		{
			name: "Properly added params to the path",
			path: "/users/{userID}",
			params: models.JSONMapStringString{
				"userID": "userID",
			},
			expectedData: "/users/userID",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			routeRepositoryMock := new(mocks.MockRouteRepository)
			routeStatusService := NewRouteStatusService(routeRepositoryMock, loggerService)
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
		expectedURL                 string
		expectedJSONData            []byte
		expectedError               error
	}
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	testsScenarios := []args{
		{
			name: "Properly prepared route for the test request",
			route: models.RouteToTest{
				RequestAuthorization: "fj349f83hf893h9834fh834",
				RequestQuery:         map[string]string{"userID": "userID"},
				Path:                 "/users/{userID}",
				RequestParams: map[string]string{
					"userID": "userID",
				},
				RequestBody: map[string]any{"appID": "appID"},
				IPAddress:   "127.0.0.1",
				Port:        "8080",
			},
			expectedAuthorizationHeader: "Bearer fj349f83hf893h9834fh834",
			expectedURL:                 "http://127.0.0.1:8080/users/userID?userID=userID",
			expectedJSONData:            []byte{0x7b, 0x22, 0x61, 0x70, 0x70, 0x49, 0x44, 0x22, 0x3a, 0x22, 0x61, 0x70, 0x70, 0x49, 0x44, 0x22, 0x7d},
			expectedError:               nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			db, err := config.NewDB(env.DBLink, "postgres")
			if err != nil {
				panic(err)
			}
			routeRepository := repository.NewRouteRepository(db.DBConnection, loggerService)
			routeStatusService := NewRouteStatusService(routeRepository, loggerService)
			authorizationHeader, url, jsonData, err := routeStatusService.prepareRouteDataForTestRequest(testScenario.route)
			assert.Equal(t, testScenario.expectedAuthorizationHeader, authorizationHeader)
			assert.Equal(t, testScenario.expectedURL, url)
			assert.Equal(t, testScenario.expectedJSONData, jsonData)
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
	env, err := config.SetConfig(tests.EnvFileLocationForServices)
	if err != nil {
		panic(err)
	}
	testsScenarios := []args{
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteBody",
			route: models.RouteToTest{
				NextRouteBody: models.JSONStringSlice{"userID"},
			},
			key:                              "userID",
			val:                              "userID",
			expectedNextRouteBody:            map[string]any{"userID": "userID"},
			expectedNextRouteParams:          map[string]string{},
			expectedRouteQuery:               map[string]string{},
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteParams",
			route: models.RouteToTest{
				NextRouteParams: models.JSONStringSlice{"userID"},
			},
			key:                              "userID",
			val:                              "userID",
			expectedNextRouteBody:            map[string]any{},
			expectedNextRouteParams:          map[string]string{"userID": "userID"},
			expectedRouteQuery:               map[string]string{},
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteQuery",
			route: models.RouteToTest{
				NextRouteQuery: models.JSONStringSlice{"userID"},
			},
			key:                              "userID",
			val:                              "userID",
			expectedNextRouteBody:            map[string]any{},
			expectedNextRouteParams:          map[string]string{},
			expectedRouteQuery:               map[string]string{"userID": "userID"},
			expectedRouteAuthorizationHeader: "",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Properly prepared route for the next request and checked routeStatus with nextRouteAuthorizationHeader",
			route: models.RouteToTest{
				NextAuthorizationHeader: "userID",
			},
			key:                              "userID",
			val:                              "userID",
			expectedNextRouteBody:            map[string]any{},
			expectedNextRouteParams:          map[string]string{},
			expectedRouteQuery:               map[string]string{},
			expectedRouteAuthorizationHeader: "userID",
			expectedRouteStatus:              "unknown",
		},
		{
			name: "Failed to assign value to the nexRouteQuery",
			route: models.RouteToTest{
				NextRouteQuery: models.JSONStringSlice{"userID"},
			},
			key:                              "userID",
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
				NextRouteParams: models.JSONStringSlice{"userID"},
			},
			key:                              "userID",
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
				NextAuthorizationHeader: "userID",
			},
			key:                              "userID",
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
			db, err := config.NewDB(env.DBLink, "postgres")
			if err != nil {
				panic(err)
			}
			routeRepository := repository.NewRouteRepository(db.DBConnection, loggerService)
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
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{}, errors.New("failed to get data from db"))
				return mRouteRepository
			},
			expectedError: errors.New("failed to get data from db"),
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
						IPAddress: "jsonplaceholder.typicode.com",
						Port:      "80",
						Path:      "/posts/{postID}",
						Method:    "GET",
						RequestParams: map[string]string{
							"postID": "1",
						},
						ResponseBody: map[string]any{
							"id":     1,
							"body":   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
							"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
							"userID": 1,
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
						IPAddress: "jsonplaceholder.typicode.com",
						Port:      "80",
						Path:      "/posts/{postID}",
						Method:    "GET",
						RequestParams: map[string]string{
							"postID": "1",
						},
						ResponseBody: map[string]any{
							"id":     1,
							"body":   "quia et suscipit suscipit recusandae consequuntur expedita et cum reprehenderit molestiae ut ut quas totam nostrum rerum est autem sunt rem eveniet architecto",
							"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
							"userID": 1,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(errors.New("failed to update working route status"))
				return mRouteRepository
			},
			expectedError: errors.New("failed to update working route status"),
		},
		{
			name: "Properly chained 2 routes",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{
					{
						IPAddress:          "192.168.0.100",
						Port:               "3040",
						Path:               "/api/v1/auth/login",
						Method:             "POST",
						ResponseStatusCode: 200,
						ParentID:           0,
						RequestBody: map[string]any{
							"email":    "adikurek121@gmail.com",
							"password": "a32lam#Fak#@ota",
						},
						ID: 72,
						ResponseBody: map[string]any{
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IlRFU1QiLCJleHAiOjE3NjE2NDM0OTYsImlkIjoxNywibmFtZSI6IlRFU1QiLCJzdXJuYW1lIjoiYWRpa3VyZWsxMjFAZ21haWwuY29tIn0.IYj7cAsrk6BeFUvHPNFOTtlG4MZqDWOXhlxRIOjJDUo",
						},
						NextAuthorizationHeader: "token",
					},
					{
						ParentID:           72,
						ID:                 73,
						IPAddress:          "192.168.0.100",
						Port:               "3040",
						Path:               "/api/v1/users/{userID}",
						Method:             "GET",
						ResponseStatusCode: 200,
						RequestParams: map[string]string{
							"userID": "17",
						},
						ResponseBody: map[string]any{
							"id":                   "17",
							"name":                 "TEST",
							"email":                "TEST",
							"surname":              "adikurek121@gmail.com",
							"password":             "$2a$10$KU6TxypBZARLVhn9ydOrDeqVGW4YnOFwuPpZF1iM/y8x4IWjleBCW",
							"createdAt":            "2025-09-30T13:11:28.430708Z",
							"updatedAt":            "2025-09-30T13:11:28.430708Z",
							"emailNotifications":   false,
							"slackNotifications":   true,
							"discordNotifications": true,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(nil)
				return mRouteRepository
			},
			expectedError: nil,
		},
		{
			name: "Wrong response status code",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{
					{
						IPAddress:          "192.168.0.100",
						Port:               "3040",
						Path:               "/api/v1/auth/login",
						Method:             "POST",
						ResponseStatusCode: 300,
						ParentID:           0,
						RequestBody: map[string]any{
							"email":    "adikurek121@gmail.com",
							"password": "a32lam#Fak#@ota",
						},
						ID: 72,
						ResponseBody: map[string]any{
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IlRFU1QiLCJleHAiOjE3NjE2NDM0OTYsImlkIjoxNywibmFtZSI6IlRFU1QiLCJzdXJuYW1lIjoiYWRpa3VyZWsxMjFAZ21haWwuY29tIn0.IYj7cAsrk6BeFUvHPNFOTtlG4MZqDWOXhlxRIOjJDUo",
						},
						NextAuthorizationHeader: "token",
					},
					{
						ParentID:           72,
						ID:                 73,
						IPAddress:          "192.168.0.100",
						Port:               "3040",
						Path:               "/api/v1/users/{userID}",
						Method:             "GET",
						ResponseStatusCode: 200,
						RequestParams: map[string]string{
							"userID": "17",
						},
						ResponseBody: map[string]any{
							"id":                   "17",
							"name":                 "TEST",
							"email":                "TEST",
							"surname":              "adikurek121@gmail.com",
							"password":             "$2a$10$KU6TxypBZARLVhn9ydOrDeqVGW4YnOFwuPpZF1iM/y8x4IWjleBCW",
							"createdAt":            "2025-09-30T13:11:28.430708Z",
							"updatedAt":            "2025-09-30T13:11:28.430708Z",
							"emailNotifications":   false,
							"slackNotifications":   true,
							"discordNotifications": true,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(nil)
				return mRouteRepository
			},
			expectedError: nil,
		},
		{
			name: "Wrong response body",
			setupMock: func() interfaces.RouteRepository {
				mRouteRepository := new(mocks.MockRouteRepository)
				mRouteRepository.On("GetWorkingRoutesToTest", mock.Anything).Return([]models.RouteToTest{
					{
						IPAddress:          "192.168.0.100",
						Port:               "3040",
						Path:               "/api/v1/auth/login",
						Method:             "POST",
						ResponseStatusCode: 200,
						ParentID:           0,
						RequestBody: map[string]any{
							"email":    "adikurek121@gmail.com",
							"password": "a32lam#Fak#@ota",
						},
						ID: 72,
						ResponseBody: map[string]any{
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IlRFU1QiLCJleHAiOjE3NjE2NDM0OTYsImlkIjoxNywibmFtZSI6IlRFU1QiLCJzdXJuYW1lIjoiYWRpa3VyZWsxMjFAZ21haWwuY29tIn0.IYj7cAsrk6BeFUvHPNFOTtlG4MZqDWOXhlxRIOjJDUo",
						},
						NextAuthorizationHeader: "token",
					},
					{
						ParentID:           72,
						ID:                 73,
						IPAddress:          "192.168.0.100",
						Port:               "3040",
						Path:               "/api/v1/users/{userID}",
						Method:             "GET",
						ResponseStatusCode: 200,
						RequestParams: map[string]string{
							"userID": "17",
						},
						ResponseBody: map[string]any{
							"id":                   "17",
							"email":                "TEST",
							"surname":              "adikurek121@gmail.com",
							"password":             "$2a$10$KU6TxypBZARLVhn9ydOrDeqVGW4YnOFwuPpZF1iM/y8x4IWjleBCW",
							"createdAt":            "2025-09-30T13:11:28.430708Z",
							"updatedAt":            "2025-09-30T13:11:28.430708Z",
							"emailNotifications":   false,
							"slackNotifications":   true,
							"discordNotifications": true,
						},
					},
				}, nil)
				mRouteRepository.On("UpdateWorkingRoutesStatuses", mock.Anything, mock.Anything).Return(nil)
				return mRouteRepository
			},
			expectedError: nil,
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
