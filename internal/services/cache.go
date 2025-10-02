package services

import (
	"context"
	"time"
)

type cacheService interface {
	SetData(ctx context.Context, key string, data string, ttl time.Duration) error
	GetData(ctx context.Context, key string) (string, error)
	ExistsData(ctx context.Context, key string) (int64, error)
	DeleteData(ctx context.Context, key string) error
}
