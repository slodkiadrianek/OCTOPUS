package utils

import (
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

func InsertionSortForRoutes[T DTO.RoutesParentId](data []T) []T {
	for i := 0; i < len(data); i++ {
		item := data[i].GetParentId()
		position := i
		for position > 0 && data[position-1].GetParentId() > item {
			data[position].SetParentId(data[position-1].GetParentId())
			position--
		}
		data[position].SetParentId(item)
	}
	return data
}
