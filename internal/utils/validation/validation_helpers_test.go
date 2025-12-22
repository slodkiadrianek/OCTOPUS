package validation

import (
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"

	z "github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/stretchr/testify/assert"
)

type ValidateSchemaTestData struct {
	Name string `json:"name"`
}

func TestValidateSchema(t *testing.T) {
	type args struct {
		name          string
		schema        z.StructSchema
		val           ValidateSchemaTestData
		expectedError z.ZogIssueMap
	}

	testsScenarios := []args{
		{
			name: "Proper data",
			schema: *z.Struct(z.Shape{
				"name": z.String(),
			}),
			val: ValidateSchemaTestData{
				Name: "test",
			},
			expectedError: z.ZogIssueMap(nil),
		},
		{
			name: "With wrong data provided",
			schema: *z.Struct(z.Shape{
				"name": z.String().Email(),
			}),
			val: ValidateSchemaTestData{
				Name: "test",
			},
			expectedError: z.ZogIssueMap{
				"$first": []*z.ZogIssue{
					{
						Code:    "email",
						Path:    "name",
						Value:   tests.Ptr("test"),
						Dtype:   "string",
						Params:  nil,
						Message: "must be a valid email",
						Err:     nil,
					},
				},
				"name": []*z.ZogIssue{
					{
						Code:    "email",
						Path:    "name",
						Value:   tests.Ptr("test"),
						Dtype:   "string",
						Params:  nil,
						Message: "must be a valid email",
						Err:     nil,
					},
				},
			},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			err := ValidateInputStruct(&testScenario.schema, &testScenario.val)
			assert.Equal(t, testScenario.expectedError, err)
		})
	}
}

func TestValidateUsersIDs(t *testing.T) {
	type args struct {
		name            string
		userID          int
		userIDFromToken int
		expectedError   error
	}
	testsScenarios := []args{
		{
			name:            "Different ids",
			userID:          2,
			userIDFromToken: 1,
			expectedError:   errors.New("provided user id's are different"),
		},
		{
			name:            "Proper data",
			userID:          1,
			userIDFromToken: 1,
			expectedError:   nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			err := ValidateUsersIDs(testScenario.userID, testScenario.userIDFromToken)
			if testScenario.expectedError != nil {
				assert.Equal(t, testScenario.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute(t *testing.T) {
	type args struct {
		name        string
		actualRoute DTO.CreateRoute
		nextRoute   DTO.CreateRoute
		result      bool
	}
	testsScenarios := []args{
		{
			name: "Proper data",
			actualRoute: DTO.CreateRoute{
				NextRouteBody: []string{"id"},
				ResponseBody: map[string]any{
					"id": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestBody: map[string]any{
					"id": "test",
				},
			},
			result: true,
		},
		{
			name: "Nested proper data",
			actualRoute: DTO.CreateRoute{
				NextRouteBody: []string{"id"},
				ResponseBody: map[string]any{
					"data": map[string]any{
						"id": "test",
					},
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestBody: map[string]any{
					"id": "test",
				},
			},
			result: true,
		},
		{
			name: "Lack of id in the response body",
			actualRoute: DTO.CreateRoute{
				NextRouteBody: []string{"id"},
				ResponseBody: map[string]any{
					"name": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestBody: map[string]any{
					"id": "test",
				},
			},
			result: false,
		},
		{
			name: "Lack of id in the next route body",
			actualRoute: DTO.CreateRoute{
				NextRouteBody: []string{"id"},
				ResponseBody: map[string]any{
					"id": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestBody: map[string]any{
					"name": "test",
				},
			},
			result: false,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res := CheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute(testScenario.actualRoute, testScenario.nextRoute)
			assert.Equal(t, testScenario.result, res)
		})
	}
}

func TestCheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute(t *testing.T) {
	type args struct {
		name        string
		actualRoute DTO.CreateRoute
		nextRoute   DTO.CreateRoute
		result      bool
	}
	testsScenarios := []args{
		{
			name: "Proper data",
			actualRoute: DTO.CreateRoute{
				NextRouteQuery: []string{"id"},
				ResponseBody: map[string]any{
					"id": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestQuery: map[string]string{"id": "test"},
			},
			result: true,
		},
		{
			name: "Nested proper data",
			actualRoute: DTO.CreateRoute{
				NextRouteQuery: []string{"id"},
				ResponseBody: map[string]any{
					"data": map[string]any{
						"id": "test",
					},
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestQuery: map[string]string{
					"id": "test",
				},
			},
			result: true,
		},
		{
			name: "Lack of id in the response body",
			actualRoute: DTO.CreateRoute{
				NextRouteQuery: []string{"id"},
				ResponseBody: map[string]any{
					"name": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestQuery: map[string]string{
					"id": "test",
				},
			},
			result: false,
		},
		{
			name: "Lack of id in the next route body",
			actualRoute: DTO.CreateRoute{
				NextRouteQuery: []string{"id"},
				ResponseBody: map[string]any{
					"id": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestQuery: map[string]string{
					"name": "test",
				},
			},
			result: false,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res := CheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute(testScenario.actualRoute, testScenario.nextRoute)
			assert.Equal(t, testScenario.result, res)
		})
	}
}

func TestCheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute(t *testing.T) {
	type args struct {
		name        string
		actualRoute DTO.CreateRoute
		nextRoute   DTO.CreateRoute
		result      bool
	}
	testsScenarios := []args{
		{
			name: "Proper data",
			actualRoute: DTO.CreateRoute{
				NextRouteParams: []string{"id"},
				ResponseBody: map[string]any{
					"id": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestParams: map[string]string{"id": "test"},
			},
			result: true,
		},
		{
			name: "Nested proper data",
			actualRoute: DTO.CreateRoute{
				NextRouteParams: []string{"id"},
				ResponseBody: map[string]any{
					"data": map[string]any{
						"id": "test",
					},
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestParams: map[string]string{
					"id": "test",
				},
			},
			result: true,
		},
		{
			name: "Lack of id in the response body",
			actualRoute: DTO.CreateRoute{
				NextRouteParams: []string{"id"},
				ResponseBody: map[string]any{
					"name": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestParams: map[string]string{
					"id": "test",
				},
			},
			result: false,
		},
		{
			name: "Lack of id in the next route body",
			actualRoute: DTO.CreateRoute{
				NextRouteParams: []string{"id"},
				ResponseBody: map[string]any{
					"id": "test",
				},
			},
			nextRoute: DTO.CreateRoute{
				RequestParams: map[string]string{
					"name": "test",
				},
			},
			result: false,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res := CheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute(testScenario.actualRoute, testScenario.nextRoute)
			assert.Equal(t, testScenario.result, res)
		})
	}
}
