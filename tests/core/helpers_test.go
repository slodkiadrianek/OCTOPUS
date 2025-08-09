package core

import (
	"errors"
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if b, ok := test.data.([]byte); ok {
				res, err := utils.UnmarshalData[map[string]string](b)
				if test.expectedError != nil {
					assert.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					assert.Equal(t, test.expectedError, err)
				}
				assert.Equal(t, test.expectedData, *res)
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
