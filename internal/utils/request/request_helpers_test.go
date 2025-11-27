package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/tests"

	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/stretchr/testify/assert"
)

func TestRemoveLastCharacterFromUrl(t *testing.T) {
	type args struct {
		name         string
		urlPath      string
		expectedData string
	}
	test := args{
		name:         "Testing remove last character from url",
		urlPath:      "/url/",
		expectedData: "/url",
	}

	t.Run(test.name, func(t *testing.T) {
		res := RemoveLastCharacterFromUrl(test.urlPath)
		assert.Equal(t, test.expectedData, res)
	})
}

func TestReadParam(t *testing.T) {
	type args struct {
		name          string
		urlPath       string
		routeKeyUrl   any
		paramToRead   string
		expectedError error
		expectedData  string
	}
	tests := []args{
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

func TestReadBody(t *testing.T) {
	type args struct {
		name          string
		bodyData      any
		expectedError error
		expectedData  any
	}
	tests := []args{
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

func TestMatchRoutes(t *testing.T) {
	type args struct {
		name         string
		routeKeyUrl  string
		urlPath      string
		expectedData bool
	}

	tests := []args{
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
	type args struct {
		name          string
		data          []string
		expectedError error
		expectedData  string
	}
	tests := []args{
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

func TestCheckRouteParams(t *testing.T) {
	type args struct {
		name           string
		actualRoute    DTO.CreateRoute
		expectedResult bool
	}

	testsScenarios := []args{
		{
			name: "Proper data provided",
			actualRoute: DTO.CreateRoute{
				RequestParams: map[string]string{"appId": "1232131", "userId": "329dfhb329"},
				Path:          "/{appId}/{userId}",
			},
			expectedResult: true,
		},
		{
			name: "Wrong path provided",
			actualRoute: DTO.CreateRoute{
				RequestParams: map[string]string{"appId": "1232131", "userId": "329dfhb329"},
				Path:          "/{appId}",
			},
			expectedResult: false,
		},
		{
			name: "Wrong  request params provided",
			actualRoute: DTO.CreateRoute{
				RequestParams: map[string]string{"userId": "329dfhb329"},
				Path:          "/{appId}",
			},
			expectedResult: false,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res := CheckRouteParams(testScenario.actualRoute)
			assert.Equal(t, testScenario.expectedResult, res)
		})
	}
}

func TestReadAllParams(t *testing.T) {
	type args struct {
		name          string
		routeKeyPath  *string
		url           string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Proper data provided",
			routeKeyPath:  tests.Ptr("/:appId/:userId"),
			url:           "/f234f3f43/3",
			expectedError: nil,
		},
		{
			name:          "Wrong routeKeyPath provided",
			routeKeyPath:  nil,
			url:           "/f234f3f43/3",
			expectedError: errors.New("failed to read context routeKeyPath, must be type string"),
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			var r *http.Request
			r = &http.Request{}
			r.URL = &url.URL{}
			r.URL.Path = testScenario.url
			if testScenario.routeKeyPath != nil {
				r = utils.SetContext(r, "routeKeyPath", *testScenario.routeKeyPath)
			} else {
				r = utils.SetContext(r, "routeKeyPath", testScenario.routeKeyPath)

			}
			res, err := ReadAllParams(r)
			fmt.Println(err)
			if testScenario.expectedError != nil {
				assert.Equal(t, err, testScenario.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NotEmpty(t, res)
				assert.Nil(t, err)
			}

		})
	}
}

func TestSendHttp(t *testing.T) {
	type args struct {
		name                string
		url                 string
		bodyFromResponse    bool
		authorizationHeader string
		method              string
		body                []byte
		expectedError       error
	}
	testsScenarios := []args{
		{
			name:                "Proper data provided with read body from response",
			url:                 "https://jsonplaceholder.typicode.com/todos/1",
			bodyFromResponse:    true,
			authorizationHeader: "",
			method:              "GET",
			body:                []byte{},
			expectedError:       nil,
		},
		{
			name:                "Proper data provided without read body from response",
			url:                 "https://jsonplaceholder.typicode.com/todos/1",
			bodyFromResponse:    false,
			authorizationHeader: "",
			method:              "GET",
			body:                []byte{},
			expectedError:       nil,
		},
		{
			name:                "Proper data provided with authorizationHeader",
			url:                 "https://jsonplaceholder.typicode.com/todos/1",
			bodyFromResponse:    false,
			authorizationHeader: "test",
			method:              "GET",
			body:                []byte{},
			expectedError:       nil,
		},
		{
			name:                "Failed to read the body",
			url:                 "https://example.com",
			bodyFromResponse:    true,
			authorizationHeader: "test",
			method:              "GET",
			body:                []byte{},
			expectedError:       errors.New("invalid character '<' looking for beginning of value"),
		},
		{
			name:                "Failed to do reuquest",
			url:                 "",
			bodyFromResponse:    true,
			authorizationHeader: "test",
			method:              "GET",
			body:                []byte{},
			expectedError:       errors.New("Get \"\": unsupported protocol scheme \"\""),
		},
		{
			name:                "Failed to create the reuquest",
			url:                 "://bad-url",
			bodyFromResponse:    true,
			authorizationHeader: "test",
			method:              "TEST",
			body:                []byte{},
			expectedError:       errors.New("parse \"://bad-url\": missing protocol scheme"),
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			statusCode, bodyFromResponse, err := SendHttp(ctx, testScenario.url, testScenario.authorizationHeader,
				testScenario.method, testScenario.body, testScenario.bodyFromResponse)
			fmt.Println(err)
			if testScenario.expectedError != nil {
				assert.Equal(t, 0, statusCode)
				assert.Equal(t, map[string]any{}, bodyFromResponse)
				assert.Equal(t, err.Error(), testScenario.expectedError.Error())
			} else {
				fmt.Println(statusCode, bodyFromResponse)
				assert.NotEqual(t, 0, statusCode)
				assert.NotEqual(t, map[string]any{}, bodyFromResponse)
				assert.Nil(t, err)
			}
		})
	}
}

func TestReadUserIdFromToken(t *testing.T) {
	type args struct {
		name          string
		id            *int
		expectedError error
	}

	testsScenarios := []args{
		{
			name:          "Proper data",
			id:            te2,
			expectedError: nil,
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			var r *http.Request
			// if testScenario.routeKeyPath != nil {
			// 	r = utils.SetContext(r, "routeKeyPath", *testScenario.routeKeyPath)
			// } else {
			// 	r = utils.SetContext(r, "routeKeyPath", testScenario.routeKeyPath)
			//
			// }
			res, err := ReadAllParams(r)
			fmt.Println(err)
			if testScenario.expectedError != nil {
				assert.Equal(t, err, testScenario.expectedError)
				assert.Equal(t, 0, res)
			} else {
				assert.Nil(t, testScenario.expectedError)
				assert.NotEqual(t, 0, res)
			}

		})
	}
}
