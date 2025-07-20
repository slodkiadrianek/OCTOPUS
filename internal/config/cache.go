package config

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	Client *redis.Client
}

func NewCacheService(cacheLink string) *CacheService {
	opt, err := redis.ParseURL(cacheLink)
	if err != nil {
		panic(err)
	}
	return &CacheService{
		Client: redis.NewClient(opt),
	}
}

func (c *CacheService) SetData(ctx context.Context, key string, data string, ttl time.Duration) error {
	err := c.Client.Set(ctx, key, string(data), ttl).Err()
	if err != nil {
		return errors.New("Failed to set value in cache")
	}
	return nil
}

func (c *CacheService) GetData(ctx context.Context, key string) (string, error) {
	res, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *CacheService) ExistsData(ctx context.Context, key string) (int64, error) {
	res, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (c *CacheService) DeleteData(ctx context.Context, key string) error {
	err := c.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
