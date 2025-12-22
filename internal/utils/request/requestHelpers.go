package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/slodkiadrianek/octopus/internal/DTO"
)

func SendHTTP(ctx context.Context, URL, authorizationHeader, method string, body []byte, readBody bool) (int,
	map[string]any, error,
) {
	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, URL, bytes.NewBuffer(body))
	if err != nil {
		return 0, map[string]any{}, err
	}
	if authorizationHeader != "" {
		req.Header.Add("Authorization", authorizationHeader)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := httpClient.Do(req)
	if err != nil {
		return 0, map[string]any{}, err
	}
	defer response.Body.Close()

	var bodyFromResponse map[string]any
	if readBody {
		err = json.NewDecoder(response.Body).Decode(&bodyFromResponse)
		fmt.Println(err)
		if err != nil {
			return 0, map[string]any{}, err
		}
	}
	return response.StatusCode, bodyFromResponse, nil
}

func ReadUserIDFromToken(r *http.Request) (int, error) {
	userID, ok := r.Context().Value("id").(int)
	if !ok || userID == 0 {
		err := errors.New("Failed to read user from context")
		return 0, err
	}
	return userID, nil
}

func ReadBody[T any](r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, errors.New("no request body provided")
	}
	var body T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func ReadQueryParam(r *http.Request, QueryName string) string {
	name := r.URL.Query().Get(QueryName)
	return name
}

func MatchRoute(routeURL, URLPath string) bool {
	splittedRouteURL := strings.Split(strings.Trim(routeURL, "/"), "/")
	splittedURLPath := strings.Split(strings.Trim(URLPath, "/"), "/")

	if len(splittedRouteURL) != len(splittedURLPath) {
		return false
	}

	for i := 0; i < len(splittedRouteURL); i++ {
		if strings.Contains(splittedRouteURL[i], ":") {
			continue
		}
		if splittedURLPath[i] != splittedRouteURL[i] {
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

func CheckRouteParams(actualRoute DTO.CreateRoute) bool {
	countParamsFromPath := 0
	splittedPath := strings.Split(actualRoute.Path, "/")
	for _, val := range splittedPath {
		leftBrace := strings.Contains(val, "{")
		rightBrace := strings.Contains(val, "}")
		if leftBrace && rightBrace {
			param := val[1 : len(val)-1]
			if _, exist := actualRoute.RequestParams[param]; !exist {
				return false
			}
			countParamsFromPath++
		}
	}
	if countParamsFromPath != len(actualRoute.RequestParams) {
		return false
	}
	return true
}

func RemoveLastCharacterFromURL(route string) string {
	if string(route[len(route)-1]) == "/" {
		route = route[:len(route)-1]
	}
	return route
}
