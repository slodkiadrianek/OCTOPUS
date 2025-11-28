package validation

import (
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
