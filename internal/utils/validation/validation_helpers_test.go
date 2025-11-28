package validation

import (
	"errors"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"testing"

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
		expectedError any
	}

	testCases := []args{
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

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateInputStruct(&test.schema, &test.val)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestValidateUsersIds(t *testing.T) {
	type args struct {
		name            string
		userId          int
		userIdFromToken int
		expectedError   error
	}
	testsScenarios := []args{
		{
			name:            "Different ids",
			userId:          2,
			userIdFromToken: 1,
			expectedError:   errors.New("Provided user id's are different"),
		},
		{
			name:            "Proper data",
			userId:          1,
			userIdFromToken: 1,
			expectedError:   nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			err := ValidateUsersIds(testScenario.userId, testScenario.userIdFromToken)
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
	}
}
