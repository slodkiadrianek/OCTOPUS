package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/models"
)

func ReadUserIdFromToken(w http.ResponseWriter, r *http.Request, logger *Logger) int {
	userId, ok := r.Context().Value("id").(int)
	if !ok || userId == 0 {
		logger.Error("Failed to read user id from context", r.URL.Path)
		err := models.NewError(500, "Server", "Internal server error")
		SetError(w, r, err)
		return 0
	}
	return userId
}

func ValidateUsersIds(w http.ResponseWriter, r *http.Request, logger *Logger, userId int, userIdToken int) {
	if userIdToken != userId {
		logger.Error("You are not allowed to do this action", map[string]any{
			"path":        r.URL.Path,
			"userIdToken": userIdToken,
		})
		err := models.NewError(500, "Server", "Internal server error")
		SetError(w, r, err)
		return
	}
}

func ReadBody[T any](r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, errors.New("no request body provided")
	}
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func ReadQueryParam(r *http.Request, QueryName string) string {
	name := r.URL.Query().Get(QueryName)
	return name
}

func MatchRoute(routeUrl, urlPath string) bool {
	splittedRouteUrl := strings.Split(strings.Trim(routeUrl, "/"), "/")
	splittedUrlPath := strings.Split(strings.Trim(urlPath, "/"), "/")

	if len(splittedRouteUrl) != len(splittedUrlPath) {
		return false
	}

	for i := 0; i < len(splittedRouteUrl); i++ {
		if strings.Contains(splittedRouteUrl[i], ":") {
			continue
		}
		if splittedUrlPath[i] != splittedRouteUrl[i] {
			return false
		}
	}
	return true
}

func ReadParam(r *http.Request, paramToRead string) (string, error) {
	path := r.URL.Path
	routeKeyPath := r.Context().Value("routeKeyPath")
	s, ok := routeKeyPath.(string)
	if !ok {
		return "", errors.New("failed to read context routeKeyPath, must be type string")
	}
	splittedPath := strings.Split(strings.Trim(path, "/"), "/")
	splittedRouteKeyPath := strings.Split(strings.Trim(s, "/"), "/")
	param := ""
	for i := 0; i < len(splittedPath); i++ {
		if strings.Contains(splittedRouteKeyPath[i], ":") && splittedRouteKeyPath[i][1:] == paramToRead {
			param = splittedPath[i]
			break
		}
	}
	if param == "" {
		return "", errors.New("The is no parameter called: " + paramToRead)
	}
	return param, nil
}

func ReadAllParams(r *http.Request) (map[string]string, error) {
	path := r.URL.Path
	routeKeyPath := r.Context().Value("routeKeyPath")
	s, ok := routeKeyPath.(string)
	if !ok {
		return nil, errors.New("failed to read context routeKeyPath, must be type string")
	}
	splittedPath := strings.Split(strings.Trim(path, "/"), "/")
	splittedRouteKeyPath := strings.Split(strings.Trim(s, "/"), "/")
	params := make(map[string]string)
	for i := 0; i < len(splittedPath); i++ {
		if strings.Contains(splittedRouteKeyPath[i], ":") {
			paramName := splittedRouteKeyPath[i][1:]
			params[paramName] = splittedPath[i]

		}
	}
	return params, nil
}

func ValidateInput(schema *z.StructSchema, val any) z.ZogIssueMap {
	errMap := schema.Validate(val)
	if errMap != nil {
		return errMap
	}
	return nil
}

func RemoveLatCharacterFromUrl(route string) string {
	if string(route[len(route)-1]) == "/" {
		route = route[:len(route)-1]
	}
	return route
}
