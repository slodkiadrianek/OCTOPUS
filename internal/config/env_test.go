package config

import (
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/tests"
	"github.com/stretchr/testify/assert"
)

func String(v string) *string {
	return &v
}

func TestReadFile(t *testing.T) {
	type args struct {
		name          string
		pathToFile    string
		expectedError error
		expectedData  map[string]string
	}
	testsScenarios := []args{
		{
			name:          "Test with proper data",
			pathToFile:    tests.TestEnvFileLocationForConfig,
			expectedError: nil,
			expectedData: map[string]string{
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
			pathToFile:    "../../.env.test2",
			expectedError: errors.New("failed to open file ../../.env.test2: open ../../.env.test2: no such file or directory"),
			expectedData:  map[string]string{},
		},
		{
			name:          "Too big file",
			pathToFile:    "../../.test.env",
			expectedError: errors.New("failed to scan a file"),
			expectedData:  map[string]string{},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res, err := readFile(testScenario.pathToFile)
			if testScenario.expectedError == nil {
				assert.Nil(t, testScenario.expectedError, err)
			} else {
				assert.Equal(t, testScenario.expectedError.Error(), err.Error())
			}
			assert.Equal(t, testScenario.expectedData, res)
		})
	}
}

func TestSetConfig(t *testing.T) {
	type args struct {
		name          string
		pathToFile    string
		expectedError error
		expectedData  Env
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			pathToFile:    ".env.test",
			expectedError: nil,
			expectedData: Env{
				Port:      "3009",
				JWTSecret: "jf3420f98234f",
				DbLink:    "postgres://adrian:zaqwerfvbgtyhn@192.168.0.100:5433/octopus?sslmode=disable",
				CacheLink: "redis://zaqwerfvbgtyhn@192.168.0.100:6379/0",
			},
		},
		{
			name:          "Wrong file provided",
			pathToFile:    ".env.test2",
			expectedError: errors.New("failed to open file .env.test2: open .env.test2: no such file or directory"),
			expectedData:  Env{},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			res, err := SetConfig(testScenario.pathToFile)

			if testScenario.expectedError == nil {
				assert.Nil(t, testScenario.expectedError, err)
			} else {
				assert.Equal(t, testScenario.expectedError.Error(), err.Error())
			}
			assert.Equal(t, testScenario.expectedData, *res)
		})
	}
}
