package utils

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContext(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	testsScenarios := []args{
		{
			key:   "testScenario",
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
			"name": "testScenario"
			}`),
			expectedError: nil,
			expectedData:  map[string]string{"name": "testScenario"},
		},
		{
			name:          "Unmarshal json with an array",
			data:          []byte(`[{"name":"testScenario"}]`),
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
				"name": "testScenario",
			},
			expectedError: nil,
			expectedData:  []byte(`{"name":"testScenario"}`),
		},
		{
			name: "Marshal object nested in array ",
			data: []map[string]string{
				{
					"name": "testScenario",
				},
			},
			expectedError: nil,
			expectedData:  []byte(`[{"name":"testScenario"}]`),
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
