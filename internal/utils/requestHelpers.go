package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/slodkiadrianek/octopus/internal/DTO"

	z "github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/models"
)

func ReadUserIdFromToken(r *http.Request) (int, error) {
	userId, ok := r.Context().Value("id").(int)
	if !ok || userId == 0 {
		err := models.NewError(500, "Server", "Internal server error")
		return 0, err
	}
	return userId, nil
}

func ValidateUsersIds(userId, userIdToken int) error {
	if userIdToken != userId {
		err := models.NewError(500, "Server", "Internal server error")
		return err
	}
	return nil
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

func ValidateInputStruct(schema *z.StructSchema, val any) z.ZogIssueMap {
	errMap := schema.Validate(val)
	if errMap != nil {
		return errMap
	}
	return nil
}

func ValidateInputSlice(schema *z.SliceSchema, val any) z.ZogIssueMap {
	errMap := schema.Validate(val)
	if errMap != nil {
		return errMap
	}
	return nil
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

func CheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute(actualRoute DTO.CreateRoute, nextRoute DTO.CreateRoute) bool {
	for _, val := range actualRoute.NextRouteBody {
		resBody := IsDataInResponseOrRequest[map[string]any](actualRoute.ResponseBody, val)
		if !resBody {
			return false
		}
		reqBody := IsDataInResponseOrRequest[map[string]any](nextRoute.RequestBody, val)
		if !reqBody {
			return false
		}
	}
	return true
}

func CheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute(actualRoute DTO.CreateRoute, nextRoute DTO.CreateRoute) bool {
	for _, val := range actualRoute.NextRouteQuery {
		resBody := IsDataInResponseOrRequest[map[string]any](actualRoute.ResponseBody, val)
		if !resBody {
			return false
		}
		reqQuery := IsDataInResponseOrRequest[map[string]string](nextRoute.RequestQuery, val)
		if !reqQuery {
			return false
		}
	}
	return true
}

func CheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute(actualRoute DTO.CreateRoute, nextRoute DTO.CreateRoute) bool {
	for _, val := range actualRoute.NextRouteParams {
		resBody := IsDataInResponseOrRequest[map[string]any](actualRoute.ResponseBody, val)
		if !resBody {
			return false
		}
		reqParams := IsDataInResponseOrRequest[map[string]string](nextRoute.RequestParams, val)
		if !reqParams {
			return false
		}
	}
	return true
}

func IsDataInResponseOrRequest[T map[string]any | map[string]string](responseBody T, data string) bool {
	switch body := any(responseBody).(type) {
	case map[string]any:
		for key, val := range body {
			if key == data {
				return true
			}
			if x, ok := val.(map[string]any); ok {
				if IsDataInResponseOrRequest(x, data) {
					return true
				}
			}
		}
	case map[string]string:
		for key := range body {
			if key == data {
				return true
			}
		}
	}
	return false
}

func RemoveLastCharacterFromUrl(route string) string {
	if string(route[len(route)-1]) == "/" {
		route = route[:len(route)-1]
	}
	return route
}
