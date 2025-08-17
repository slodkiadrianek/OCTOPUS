package auth

import (
	"testing"

	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/tests/api"
	"github.com/stretchr/testify/assert"
)

var lastInsertUserId string

type ProperTestRegisterUserCase struct {
	name         string
	data         map[string]string
	expectedCode int
	expectedData map[string]string
}

type ErrorTestRegisterUserCase struct {
	name         string
	data         map[string]string
	expectedCode int
	expectedData []map[string]string
}

func TestRegisterUser(t *testing.T) {
	properTestCases := []ProperTestRegisterUserCase{
		{
			name: "Test with proper data",
			data: map[string]string{
				"name":     "TEST",
				"surname":  "TEST",
				"email":    "adiku@gmail.com",
				"password": "a32lam#Fak#@ota",
			},
			expectedCode: 201,
			expectedData: map[string]string{},
		},
	}
	errorTestCases := []ErrorTestRegisterUserCase{
		{
			name: "Not proper email address",
			data: map[string]string{
				"name":     "TEST",
				"surname":  "TEST",
				"email":    "adigmail.com",
				"password": "a32lam#Fak#@ota",
			},
			expectedCode: 422,
			expectedData: []map[string]string{
				{
					"Code":    "email",
					"Dtype":   "string",
					"Err":     "",
					"Message": "must be a valid email",
					"Params":  "",
					"Path":    "email",
					"Value":   "adigmail.com",
				},
			},
		},
	}
	for _, properTestCase := range properTestCases {
		CreateTestForRegisterUser[map[string]string](properTestCase.name, properTestCase.data, properTestCase.expectedCode, properTestCase.expectedData, t)
	}
	for _, errorTestCase := range errorTestCases {
		CreateTestForRegisterUser[[]map[string]string](errorTestCase.name, errorTestCase.data, errorTestCase.expectedCode, errorTestCase.expectedData, t)
	}
}

func CreateTestForRegisterUser[T any](name string, data map[string]string, expectedCode int, expectedData T, t *testing.T) {
	t.Run(name, func(t *testing.T) {
		router := api.SetupRouter()
		recorder := api.PerformTestRequest(router, "POST", "/api/v1/auth/register", data)
		assert.Equal(t, expectedCode, recorder.Code)
		resBody, err := utils.UnmarshalData[T](recorder.Body.Bytes())
		if err != nil {
			panic(err)
		}
		assert.Equal(t, &expectedData, resBody)
		if expectedCode == 201 {
			err = api.DbInterface.DbConn.QueryRow(`SELECT currval('users_id_seq')`).Scan(&lastInsertUserId)
			if err != nil {
				panic(err)
			}
			CleanDbUsersTable()
		}
	})
}

func CleanDbUsersTable() {
	_, err := api.DbInterface.DbConn.Exec("DELETE FROM users WHERE id =" + lastInsertUserId)
	if err != nil {
		panic(err)
	}
}
