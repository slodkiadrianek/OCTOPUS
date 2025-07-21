package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/slodkiadrianek/octopus/pkg/errors"
)

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
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}
}

func ReadBody[T any](w http.ResponseWriter, r *http.Request, model T) T {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		SendResponse(w, 500, errors.Err_body_res)
	}
	return body
}

func readQueryParam(r *http.Request, QueryName string) string {
	name := r.URL.Query().Get(QueryName)
	return name
}

func readParams(r *http.Request, paramsToRead []string) map[string]string {
	path := r.URL.Path
	splitPath := strings.Split(path, "/")
	params := make(map[string]string)
	for i := 0; i < len(splitPath); i++ {
		for _, val := range paramsToRead {
			if splitPath[i] == val {
				if i+1 < len(splitPath) {
					params[val] = splitPath[i+1]
				}
			}
		}
	}
	return params
}
