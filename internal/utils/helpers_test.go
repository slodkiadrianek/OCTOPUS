package utils

import (
	"errors"
	"net/http"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/stretchr/testify/assert"
)

func TestSetContext(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	testsScenarios := []args{
		{
			key:   "test",
			value: "testValue",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run("Testing set context", func(t *testing.T) {
			var r http.Request
			res := SetContext(&r, testScenario.key, testScenario.value)
			assert.Equal(t, testScenario.value, res.Context().Value(testScenario.key))
		})
	}
}

func TestUnmarshalData(t *testing.T) {
	type args struct {
		name          string
		data          any
		expectedError error
		expectedData  any
	}
	testsScenarios := []args{
		{
			name: "Unmarshal simple object",
			data: []byte(`{
			"name": "test"
			}`),
			expectedError: nil,
			expectedData:  map[string]string{"name": "test"},
		},
		{
			name:          "Unmarshal json with an array",
			data:          []byte(`[{"name":"test"}]`),
			expectedError: errors.New("json: cannot unmarshal array into Go value of type map[string]string"),
			expectedData:  nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			if b, ok := testScenario.data.([]byte); ok {
				res, err := UnmarshalData[map[string]string](b)
				if testScenario.expectedError != nil {
					assert.Equal(t, testScenario.expectedError.Error(), err.Error())
					assert.Nil(t, testScenario.expectedData, res)
				} else {
					assert.Equal(t, testScenario.expectedError, err)
					assert.Equal(t, testScenario.expectedData, *res)
				}
			}
		})
	}
}

func TestMarshalData(t *testing.T) {
	type args struct {
		name          string
		data          any
		expectedError error
		expectedData  any
	}
	testsScenarios := []args{
		{
			name: "Marshal simple object",
			data: map[string]string{
				"name": "test",
			},
			expectedError: nil,
			expectedData:  []byte(`{"name":"test"}`),
		},
		{
			name: "Marshal object nested in array ",
			data: []map[string]string{
				{
					"name": "test",
				},
			},
			expectedError: nil,
			expectedData:  []byte(`[{"name":"test"}]`),
		},
		{
			name:          "Create error in Marshal data fn",
			data:          func() {},
			expectedError: errors.New("json: unsupported type: func()"),
			expectedData:  []byte{},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res, err := MarshalData(testScenario.data)
			if testScenario.expectedError != nil {
				assert.Equal(t, testScenario.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, testScenario.expectedError, err)
			}
			assert.Equal(t, testScenario.expectedData, res)
		})
	}
}

func TestGenerateID(t *testing.T) {
	type args struct {
		name          string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			generatedID, err := GenerateID()
			if testScenario.expectedError != nil {
				assert.Equal(t, testScenario.expectedError.Error(), err.Error())
				assert.Equal(t, "", generatedID)
			} else {
				assert.NotEqual(t, "", generatedID)
				assert.Nil(t, err)
			}
		})
	}
}

func TestInsertionSortForRoutes(t *testing.T) {
	type args struct {
		name            string
		routesData      []DTO.RoutesParentID
		sortedRouteData []DTO.RoutesParentID
		expectedError   error
	}
	testsScenarios := []args{
		{
			name: "Proper data",
			routesData: []DTO.RoutesParentID{
				DTO.NewRouteInfo("", "", 2),
				DTO.NewRouteInfo("", "", 1),
			},
			sortedRouteData: []DTO.RoutesParentID{
				DTO.NewRouteInfo("", "", 1),
				DTO.NewRouteInfo("", "", 2),
			},
			expectedError: nil,
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			sortedRouteData := InsertionSortForRoutes(testScenario.routesData)
			assert.Equal(t, testScenario.sortedRouteData, sortedRouteData)
		})
	}
}
