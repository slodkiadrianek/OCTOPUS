package config

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCacheService(t *testing.T) {
	type args struct {
		name          string
		cacheLink     string
		expectedError *string
		expectedData  bool
	}
	testCases := []args{
		{
			name:          "Proper data",
			cacheLink:     "redis://zaqwerfvbgtyhn@192.168.0.100:6379/0",
			expectedError: nil,
			expectedData:  true,
		},
		{
			name:          "Wrong connection link",
			cacheLink:     "://aqwerfvbgtyhn@192.168.0.100:6379/0",
			expectedError: String(`parse "://aqwerfvbgtyhn@192.168.0.100:6379/0": missing protocol scheme`),
			expectedData:  false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := NewCacheService(testCase.cacheLink)
			if testCase.expectedError == nil {
				assert.Nil(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, *testCase.expectedError, err.Error())
			}
			if res == nil || *res != (CacheService{}) {
				assert.Equal(t, testCase.expectedData, true)
			} else {
				assert.Equal(t, testCase.expectedData, false)
			}
		})
	}
}

func TestCacheService_SetData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError *string
	}

	testCases := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
		},
		{
			name:          "Wrong  data",
			key:           strings.Repeat("x", 600*1024*1024),
			expectedError: String("write: connection reset by peer"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")
			ctx := context.Background()
			ttl := 20 * time.Millisecond
			err := serviceClient.SetData(ctx, testCase.key, "hi", ttl)
			if testCase.expectedError == nil {
				assert.Nil(t, err) // TTL truncation does NOT count as an error
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), *testCase.expectedError)
			}
		})
	}
}

func TestCacheService_GetData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError *string
		expectedData  string
	}
	testCases := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
			expectedData:  "h1",
		},
		{
			name:          "Wrong  data",
			key:           "",
			expectedError: String(`redis: nil`),
			expectedData:  "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")

			ctx := context.Background()
			ttl := 20 * time.Millisecond
			if testCase.expectedError == nil {
				_ = serviceClient.SetData(ctx, testCase.key, "h1", ttl)
			}
			res, err := serviceClient.GetData(ctx, testCase.key)
			if testCase.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), *testCase.expectedError)
			}
			assert.Equal(t, testCase.expectedData, res)
		})
	}
}

func TestCacheService_ExistsData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError *string
		expectedData  int64
	}
	testCases := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
			expectedData:  1,
		},
		{
			name:          "Wrong  data",
			key:           strings.Repeat("x", 600*1024*1024),
			expectedError: String("write: connection reset by peer"),
			expectedData:  0,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")

			ctx := context.Background()
			ttl := 20 * time.Millisecond
			if testCase.expectedError == nil {
				_ = serviceClient.SetData(ctx, testCase.key, "h1", ttl)
			}
			res, err := serviceClient.ExistsData(ctx, testCase.key)
			if testCase.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), *testCase.expectedError)
			}
			assert.Equal(t, testCase.expectedData, res)
		})
	}
}

func TestCacheService_DeleteData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError *string
	}
	testCases := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
		},
		{
			name:          "Wrong  data",
			key:           strings.Repeat("x", 600*1024*1024),
			expectedError: String("write: connection reset by peer"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")

			ctx := context.Background()
			ttl := 20 * time.Millisecond
			if testCase.expectedError == nil {
				_ = serviceClient.SetData(ctx, testCase.key, "h1", ttl)
			}
			err := serviceClient.DeleteData(ctx, testCase.key)
			if testCase.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), *testCase.expectedError)
			}
		})
	}
}
