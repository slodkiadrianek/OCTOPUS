package mocks

import (
	"context"
	"errors"
	"time"
)

type MockCacheService struct{}

func (m *MockCacheService) SetData(ctx context.Context, key string, data string, ttl time.Duration) error {
	return errors.New("failed to set data")
}

func (m *MockCacheService) GetData(ctx context.Context, key string) (string, error) {
	return "", errors.New("failed to get data")
}
func (m *MockCacheService) ExistsData(ctx context.Context, key string) (int64, error) {
	return 0, errors.New("cache unavailable")
}
func (m *MockCacheService) DeleteData(ctx context.Context, key string) error {
	return errors.New("cache unavailable")
}

func NewMockCacheService() *MockCacheService {
	return &MockCacheService{}
}
