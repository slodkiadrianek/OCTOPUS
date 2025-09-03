package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	z "github.com/Oudwins/zog"
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

type testReadParamData struct {
	name          string
	urlPath       string
	routeKeyUrl   any
	paramToRead   string
	expectedError error
	expectedData  string
}

type testMatchRoutesData struct {
	name         string
	routeKeyUrl  string
	urlPath      string
	expectedData bool
}

type testRemoveLastCharacterFromUrlData struct {
	name         string
	urlPath      string
	expectedData string
}
type testReadBodyData struct {
	name          string
	bodyData      any
	expectedError error
	expectedData  any
}

type testValidateSchemaData struct {
	name          string
	schema        z.StructSchema
	val           ValidateSchemaTestData
	expectedError any
}

type ValidateSchemaTestData struct {
	Name string `json:"name"`
}

func ptr(s string) *string {
	return &s
}

func TestSetContext(t *testing.T){
	type args struct{
		key string
		value any
	}
	tests := []args{
		{
			key: "test",
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

func TestValidateSchema(t *testing.T) {
	testCases := []testValidateSchemaData{
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
						Value:   ptr("test"),
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
						Value:   ptr("test"),
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
			err := ValidateInput(&test.schema, &test.val)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestReadBody(t *testing.T) {
	tests := []testReadBodyData{
		{
			name:          "Test with proper data",
			bodyData:      `{"name":"test"}`,
			expectedError: nil,
			expectedData:  map[string]string{"name": "test"},
		},
		{
			name:          "Test with malformed json",
			bodyData:      `{this is invalid json}`,
			expectedError: errors.New("invalid character 't' looking for beginning of object key string"),
			expectedData:  nil,
		},
		{
			name:          "Test without body",
			bodyData:      nil,
			expectedError: errors.New("no request body provided"),
			expectedData:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var r http.Request
			if test.bodyData == nil {
				r.Body = nil
			} else {
				s, ok := test.bodyData.(string)
				if !ok {
					panic(ok)
				}
				r.Body = io.NopCloser(bytes.NewBufferString(s))
			}
			res, err := ReadBody[map[string]string](&r)
			if test.expectedError != nil {
				fmt.Println(err.Error())
				assert.Equal(t, test.expectedError.Error(), err.Error())
				assert.Equal(t, test.expectedData, nil)
			} else {
				assert.Equal(t, test.expectedError, nil)
				assert.Equal(t, test.expectedData, *res)
			}
		})
	}
}

func TestRemoveLastCharacterFromUrl(t *testing.T) {
	test := testRemoveLastCharacterFromUrlData{
		name:         "Testing remove last character from url",
		urlPath:      "/url/",
		expectedData: "/url",
	}

	t.Run(test.name, func(t *testing.T) {
		res := RemoveLatCharacterFromUrl(test.urlPath)
		assert.Equal(t, test.expectedData, res)
	})
}

func TestReadParam(t *testing.T) {
	tests := []testReadParamData{
		{
			name:          "Proper urlPath and expectedData with 1 param in path",
			urlPath:       "/users/1",
			routeKeyUrl:   "/users/:userId",
			paramToRead:   "userId",
			expectedError: nil,
			expectedData:  "1",
		},
		{
			name:          "Proper urlPath and expectedData with 2 params in path",
			urlPath:       "/users/1/posts/1",
			routeKeyUrl:   "/users/:userId/posts/:postId",
			paramToRead:   "postId",
			expectedError: nil,
			expectedData:  "1",
		},
		{
			name:          "lack off the requested param",
			urlPath:       "/users/1",
			routeKeyUrl:   "/users/:userId",
			paramToRead:   "postId",
			expectedError: errors.New("The is no parameter called: postId"),
			expectedData:  "",
		},
		{
			name:          "Wrong type for value stored in context",
			urlPath:       "/users/1",
			routeKeyUrl:   1,
			paramToRead:   "postId",
			expectedError: errors.New("failed to read context routeKeyPath, must be type string"),
			expectedData:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var r http.Request
			r.URL = &url.URL{}
			r.URL.Path = test.urlPath
			ctx := context.WithValue(r.Context(), "routeKeyPath", test.routeKeyUrl)
			r = *r.WithContext(ctx)
			res, err := ReadParam(&r, test.paramToRead)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, nil)
			}
			assert.Equal(t, test.expectedData, res)
		})
	}
}

func TestMatchRoutes(t *testing.T) {
	tests := []testMatchRoutesData{
		{
			name:         "Test same urls",
			routeKeyUrl:  "/url/v1/v1",
			urlPath:      "/url/v1/v1",
			expectedData: true,
		},
		{
			name:         "Test different urls but with the same length",
			routeKeyUrl:  "/ur2",
			urlPath:      "/ur1",
			expectedData: false,
		},
		{
			name:         "Test urls with different lengths",
			routeKeyUrl:  "/url",
			urlPath:      "/url/12232",
			expectedData: false,
		},
		{
			name:         "Test urls with different lengths",
			routeKeyUrl:  "/url1",
			urlPath:      "/url1",
			expectedData: true,
		},
		{
			name:         "Test urls with parameters included in path",
			routeKeyUrl:  "/url1/:id/123",
			urlPath:      "/url1/1/123",
			expectedData: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := MatchRoute(test.routeKeyUrl, test.urlPath)
			assert.Equal(t, test.expectedData, res)
		})
	}
}

func TestReadQueryParam(t *testing.T) {
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
			res := ReadQueryParam(&r, test.data[0])
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
