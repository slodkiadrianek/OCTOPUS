package auth

import (
	"fmt"
	"testing"

	"github.com/slodkiadrianek/octopus/tests/api"
	"github.com/stretchr/testify/assert"
)


func TestRegisterUser(t *testing.T){
	router := api.SetupRouter()
	recorder := api.PerformTestRequest(router, "POST", "/api/v1/auth/register", map[string]string{
		"name":"TEST",
		"surname":"TEST",
		"email":"adikurek1221@gmail.com",
		"password": "a32lam#Fak#@ota",
	})
	fmt.Println(recorder)
	assert.Equal(t, 201, recorder.Code, "Expected status code 201 Created")	
}