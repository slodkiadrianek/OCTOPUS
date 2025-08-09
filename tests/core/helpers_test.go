package core

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	name          string
	data          any
	expectedError error
	expectedData  any
}

type testReadQueryParamData struct {
	name          string
	data          []string
	expectedError error
	expectedData  string
}

type testReadParamsData struct {
	name          string
	data          []string
	expectedError error
	expectedData  map[string]string
}

func TestReadQeuryParam(t *testing.T) {
	tests := []testReadQueryParamData{
		{
			name:          "Read query param properly",
			data:          []string{"name", "test"},
			expectedError: nil,
			expectedData:  "test",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var r http.Request
			r.URL = &url.URL{}
			q := r.URL.Query()
			q.Add(test.data[0], test.data[1])
			r.URL.RawQuery = q.Encode()
			res := utils.ReadQueryParam(&r, test.data[0])
			assert.Equal(t, test.expectedError, nil)
			assert.Equal(t, test.expectedData, res)
		})
	}
}

func TestUnmarshalData(t *testing.T) {
	tests := []testData{
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if b, ok := test.data.([]byte); ok {
				res, err := utils.UnmarshalData[map[string]string](b)
				if test.expectedError != nil {
					assert.Equal(t, test.expectedError.Error(), err.Error())
					assert.Nil(t, test.expectedData, res)
				} else {
					assert.Equal(t, test.expectedError, err)
					assert.Equal(t, test.expectedData, *res)
				}
			}
		})
	}
}

func TestMarshalData(t *testing.T) {
	tests := []testData{
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := utils.MarshalData(test.data)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, err)
			}
			assert.Equal(t, test.expectedData, res)
		})
	}
}
