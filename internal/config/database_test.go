package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	type args struct {
		name          string
		driver        string
		databaseLink  string
		expectedError error
		expectedData  bool
	}

	testCases := []args{
		{
			name:          "Proper link",
			driver:        "postgres",
			databaseLink:  "postgres://adrian:zaqwerfvbgtyhn@192.168.0.100:5433/octopus?sslmode=disable",
			expectedError: nil,
			expectedData:  true,
		}, {
			name:          "Wronk connection link",
			driver:        "postgres",
			databaseLink:  "invalid_connection_string",
			expectedError: errors.New(`missing "=" after "invalid_connection_string" in connection info string"`),
			expectedData:  false,
		}, {
			name:          "Wronk driver ",
			driver:        "unk",
			databaseLink:  "invalid_connection_string",
			expectedError: errors.New(`sql: unknown driver "unk" (forgotten import?)`),
			expectedData:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := NewDB(testCase.databaseLink, testCase.driver)
			if testCase.expectedError == nil {
				assert.Nil(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			}
			if res == nil || *res != (DB{}) {
				assert.Equal(t, testCase.expectedData, true)
			} else {
				assert.Equal(t, testCase.expectedData, false)
			}
		})
	}
}
