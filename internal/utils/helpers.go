package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
)

func GenerateID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func SetContext(r *http.Request, key any, data any) *http.Request {
	ctx := context.WithValue(r.Context(), key, data)
	return r.WithContext(ctx)
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

func InsertionSortForRoutes[T DTO.RoutesParentID](data []T) []T {
	for i := 1; i < len(data); i++ {
		key := data[i]
		j := i - 1
		for j >= 0 && data[j].GetParentID() > key.GetParentID() {
			data[j+1] = data[j]
			j--
		}
		data[j+1] = key
	}
	return data
}

func DoHttpRequest(ctx context.Context, url, authorizationHeader, method string, body []byte, loggerService Logger) (int,
	map[string]any, error,
) {
	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		loggerService.Error("Failed to create webhook request", err)
		return 0, map[string]any{}, err
	}

	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := httpClient.Do(req)
	if err != nil {
		loggerService.Error("Failed to send webhook request", err)
		return 0, map[string]any{}, err
	}
	defer response.Body.Close()

	var bodyFromResponse map[string]any

	err = json.NewDecoder(response.Body).Decode(&bodyFromResponse)
	if err != nil {
		return 0, map[string]any{}, err
	}

	return response.StatusCode, bodyFromResponse, nil
}
