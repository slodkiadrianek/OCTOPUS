package servicesApp

import (
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/tests"
	"testing"
)

func TestSortRoutesToTest(t *testing.T) {
	type args struct{
		name string
		routeToTest []models.RouteToTest
		expectedData map[string][]models.RouteToTest
	}
	testsScenarios := []args{
		{
			name: "Properly sorted routes",
			routeToTest: []models.RouteToTest{
				{
						ParentID: 0,
						Name: "First Route",
						AppId: "1",
				},
				{
					ParentID: 1,
					Name: "Second Route",
					AppId: "2",
				},
				{
					ParentID: 0,
					Name: "Third Route",
					AppId: "3",
				},
			},
			expectedData: map[string][]models.RouteToTest{
				"First Route1": {
					{
						ParentID: 0,
						Name: "First Route",
						AppId: "1",
					},
					{
						ParentID: 1,
						Name: "Second Route",
						AppId: "2",
					},
				},
				"Third Route3": {
					{
						ParentID: 0,
						Name: "Third Route",
						AppId: "3",
					},
				},
			},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			loggerService := tests.CreateLogger()
			db := tests.
			routeRepository := repository.NewRouteRepository(db, loggerService)
			routeStatusService := NewRouteStatusService()
			sortedData := sortR
		})
	}
}
