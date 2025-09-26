package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/models"
)

func GenerateID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

type contextKey string

const ErrorKey contextKey = "Error"

func StripANSI(str string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mGKH]`)
	return ansiRegex.ReplaceAllString(str, "")
}
func SetContext(r *http.Request, key any, data any) *http.Request {
	ctx := context.WithValue(r.Context(), key, data)
	return r.WithContext(ctx)
}

func SetError(w http.ResponseWriter, r *http.Request, err error) {
	errBucket, ok := r.Context().Value("ErrorBucket").(*models.ErrorBucket)
	if ok {
		errBucket.Err = err
		return
	}
	SendResponse(w, 500, map[string]string{
		"errorCategory":    "Server",
		"errorDescription": "Internal server error",
	})
}

func MarshalData(data any) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	return dataBytes, nil
}

func UnmarshalData[T any](dataBytes []byte) (*T, error) {
	var data *T
	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SendResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == 204 {
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
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
