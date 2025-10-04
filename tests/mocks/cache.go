package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) SetData(ctx context.Context, key string, data string, ttl time.Duration) error {
	args := m.Called(ctx, key, data, ttl)
	return args.Error(0)
}

func (m *MockCacheService) GetData(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) ExistsData(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) DeleteData(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
