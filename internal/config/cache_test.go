package config

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCacheService(t *testing.T) {
	type args struct {
		name          string
		cacheLink     string
		expectedError error
		expectedData  bool
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			cacheLink:     "redis://zaqwerfvbgtyhn@192.168.0.100:6379/0",
			expectedError: nil,
			expectedData:  true,
		},
		{
			name:          "Wrong connection link",
			cacheLink:     "://aqwerfvbgtyhn@192.168.0.100:6379/0",
			expectedError: errors.New(`parse "://aqwerfvbgtyhn@192.168.0.100:6379/0": missing protocol scheme`),
			expectedData:  false,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res, err := NewCacheService(testScenario.cacheLink)

			if testScenario.expectedError == nil {
				assert.Nil(t, testScenario.expectedError, err)
			} else {
				assert.Equal(t, testScenario.expectedError.Error(), err.Error())
			}

			if res == nil || *res != (CacheService{}) {
				assert.Equal(t, testScenario.expectedData, true)
			} else {
				assert.Equal(t, testScenario.expectedData, false)
			}
		})
	}
}

func TestCacheService_SetData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError error
	}

	testsScenarios := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
		},
		{
			name:          "Wrong  data",
			key:           strings.Repeat("x", 600*1024*1024),
			expectedError: errors.New("write: connection reset by peer"),
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")
			ctx := context.Background()
			ttl := 20 * time.Millisecond

			err := serviceClient.SetData(ctx, testScenario.key, "hi", ttl)
			if testScenario.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestCacheService_GetData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError error
		expectedData  string
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
			expectedData:  "h1",
		},
		{
			name:          "Wrong  data",
			key:           "",
			expectedError: errors.New(`redis: nil`),
			expectedData:  "",
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")

			ctx := context.Background()
			ttl := 20 * time.Millisecond
			if testScenario.expectedError == nil {
				_ = serviceClient.SetData(ctx, testScenario.key, "h1", ttl)
			}
			res, err := serviceClient.GetData(ctx, testScenario.key)
			if testScenario.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			assert.Equal(t, testScenario.expectedData, res)
		})
	}
}

func TestCacheService_ExistsData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError error
		expectedData  int64
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
			expectedData:  1,
		},
		{
			name:          "Wrong  data",
			key:           strings.Repeat("x", 600*1024*1024),
			expectedError: errors.New("write: connection reset by peer"),
			expectedData:  0,
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")

			ctx := context.Background()
			ttl := 20 * time.Millisecond
			if testScenario.expectedError == nil {
				_ = serviceClient.SetData(ctx, testScenario.key, "h1", ttl)
			}
			res, err := serviceClient.ExistsData(ctx, testScenario.key)
			if testScenario.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
			assert.Equal(t, testScenario.expectedData, res)
		})
	}
}

func TestCacheService_DeleteData(t *testing.T) {
	type args struct {
		name          string
		key           string
		expectedError error
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			key:           "test",
			expectedError: nil,
		},
		{
			name:          "Wrong  data",
			key:           strings.Repeat("x", 600*1024*1024),
			expectedError: errors.New("write: connection reset by peer"),
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			serviceClient, _ := NewCacheService("redis://:zaqwerfvbgtyhn@192.168.0.100:6379/0")

			ctx := context.Background()
			ttl := 20 * time.Millisecond
			if testScenario.expectedError == nil {
				_ = serviceClient.SetData(ctx, testScenario.key, "h1", ttl)
			}
			err := serviceClient.DeleteData(ctx, testScenario.key)
			if testScenario.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}
