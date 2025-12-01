package servicesApp

import (
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/stretchr/testify/assert"

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

func TestPrepareRouteDataForTestRequest(t *testing.T) {
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
