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
	tests := []args{
		{
			key:   "test",
			value: "testValue",
		},
	}
	for _, test := range tests {
		t.Run("Testing set context", func(t *testing.T) {
			var r http.Request
			res := SetContext(&r, test.key, test.value)
			assert.Equal(t, test.value, res.Context().Value(test.key))
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
	tests := []args{
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
				res, err := UnmarshalData[map[string]string](b)
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
	type args struct {
		name          string
		data          any
		expectedError error
		expectedData  any
	}
	tests := []args{
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
			res, err := MarshalData(test.data)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, err)
			}
			assert.Equal(t, test.expectedData, res)
		})
	}
}
