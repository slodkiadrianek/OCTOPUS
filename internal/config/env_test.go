package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestReadFileCase struct {
	name          string
	pathToFile    string
	expectedError *string
	expectedData  *map[string]string
}

func TestReadFile(t *testing.T) {
	testCases := []TestReadFileCase{
		{
			name:          "Test with proper data",
			pathToFile:    "../../.env.test",
			expectedError: nil,
			expectedData: &map[string]string{
				"Port":      "3009",
				"JWTSecret": "jf3420f98234f",
				"DbLink":    "postgres://adrian:zaqwerfvbgtyhn@192.168.0.100:5433/octopus?sslmode=disable",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := ReadFile(testCase.pathToFile)
			if testCase.expectedError == nil {
				assert.Nil(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, testCase.expectedError, err)
			}
			assert.Equal(t, *testCase.expectedData, res)
		})
	}
}

func TestSetConfig(t *testing.T) {
}

