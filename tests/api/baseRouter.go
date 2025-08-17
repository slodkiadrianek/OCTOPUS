package api

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"

	"github.com/slodkiadrianek/octopus/internal/api"
	"github.com/slodkiadrianek/octopus/internal/api/routes"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/controllers"
	"github.com/slodkiadrianek/octopus/internal/repository"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type Db struct {
	DbConn *sql.DB
}

var DbInterface Db

type TestRouter struct {
	Router *routes.Router
}

func NewTestDependencies() *TestRouter {
	loggerService := logger.NewLogger("../../../logs", "02.01.2006")
	cfg := config.SetConfig("../../../.env")
	db := config.NewDb(cfg.DbLink)
	DbInterface.DbConn = db.DbConnection
	userRepository := repository.NewUserRepository(db.DbConnection, loggerService)
	userService := services.NewUserService(loggerService, userRepository)
	userController := controllers.NewUserController(userService)
	dependenciesConfig := api.NewDependencyConfig(cfg.Port, userController)
	server := api.NewServer(dependenciesConfig)
	server.SetupMiddleware()
	server.SetupRoutes()
	return &TestRouter{
		Router: server.Router,
	}
}

func SetupRouter() *routes.Router {
	return NewTestDependencies().Router
}

func CreateNewTestRequest(method, path string, body map[string]string) (*http.Request, error) {
	jsonBody, err := utils.MarshalData(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func PerformTestRequest(router *routes.Router, method, url string, body map[string]string) *httptest.ResponseRecorder {
	req, _ := CreateNewTestRequest(method, url, body)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}
