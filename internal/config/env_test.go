package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func String(v string) *string {
	return &v
}

func TestReadFile(t *testing.T) {
	type args struct {
		name          string
		pathToFile    string
		expectedError *string
		expectedData  *map[string]string
	}
	testCases := []args{
		{
			name:          "Test with proper data",
			pathToFile:    "../../.env.test",
			expectedError: nil,
			expectedData: &map[string]string{
				"Port":         "3009",
				"JWTSecret":    "jf3420f98234f",
				"DbLink":       "postgres://adrian:zaqwerfvbgtyhn@192.168.0.100:5433/octopus?sslmode=disable",
				"CacheLink":    "redis://zaqwerfvbgtyhn@192.168.0.100:6379/0",
				"EmailService": "fj3402f",
				"EmailUser":    "32jf023hnf2",
				"EmailFrom":    "32rj9230hf",
				"EmailPass":    "390f32hjf",
			},
		},
		{
			name:          "Read file which does not exist",
			pathToFile:    "../../../.env.test",
			expectedError: String("failed to open a file"),
			expectedData:  &map[string]string{},
		},
		{
			name:          "Too big file",
			pathToFile:    "../../test.env",
			expectedError: String("failed to scan a file"),
			expectedData:  &map[string]string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := ReadFile(testCase.pathToFile)
			if testCase.expectedError == nil {
				assert.Nil(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, *testCase.expectedError, err.Error())
			}
			assert.Equal(t, *testCase.expectedData, res)
		})
	}
}

func TestSetConfig(t *testing.T) {
	type args struct {
		name          string
		pathToFile    string
		expectedError *string
		expectedData  Env
	}
	testCases := []args{
		{
			name:          "Proper data",
			pathToFile:    "../../.env.test",
			expectedError: nil,
			expectedData: Env{
				Port:         "3009",
				JWTSecret:    "jf3420f98234f",
				DbLink:       "postgres://adrian:zaqwerfvbgtyhn@192.168.0.100:5433/octopus?sslmode=disable",
				CacheLink:    "redis://zaqwerfvbgtyhn@192.168.0.100:6379/0",
				EmailService: "fj3402f",
				EmailUser:    "32jf023hnf2",
				EmailFrom:    "32rj9230hf",
				EmailPass:    "390f32hjf",
			},
		},
		{
			name:          "Wrong file provided",
			pathToFile:    "../../../.env.test",
			expectedError: String("failed to open a file"),
			expectedData:  Env{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := SetConfig(testCase.pathToFile)

			if testCase.expectedError == nil {
				assert.Nil(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, *testCase.expectedError, err.Error())
			}
			assert.Equal(t, testCase.expectedData, *res)
		})
	}
}
