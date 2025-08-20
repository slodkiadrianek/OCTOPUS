package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDb(t *testing.T) {
	type args struct {
		name          string
		driver        string
		databaseLink  string
		expectedError *string
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
			expectedError: String(`missing "=" after "invalid_connection_string" in connection info string"`),
			expectedData:  false,
		}, {
			name:          "Wronk driver ",
			driver:        "unk",
			databaseLink:  "invalid_connection_string",
			expectedError: String(`sql: unknown driver "unk" (forgotten import?)`),
			expectedData:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := NewDb(testCase.databaseLink, testCase.driver)
			if testCase.expectedError == nil {
				assert.Nil(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, *testCase.expectedError, err.Error())
			}
			if res == nil || *res != (Db{}) {
				assert.Equal(t, testCase.expectedData, true)
			} else {
				assert.Equal(t, testCase.expectedData, false)
			}
		})
	}
}
