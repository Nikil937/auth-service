package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRefreshRepository struct {
	client *redis.Client
}

var _ RefreshRepository = (*RedisRefreshRepository)(nil)

func NewRedisRefreshRepository(client *redis.Client) *RedisRefreshRepository {
	return &RedisRefreshRepository{
		client: client,
	}
}

func (r *RedisRefreshRepository) Save(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error {
	err := r.client.Set(ctx, token, userID.String(), ttl).Err()
	if err != nil {
		return fmt.Errorf("save refresh session: %w", err)
	}

	return nil
}

func (r *RedisRefreshRepository) Get(ctx context.Context, token string) (uuid.UUID, error) {
	val, err := r.client.Get(ctx, token).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("get refresh session: %w", err)
	}

	userID, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse user id: %w", err)
	}

	return userID, nil
}

func (r *RedisRefreshRepository) Delete(ctx context.Context, token string) error {
	err := r.client.Del(ctx, token).Err()
	if err != nil {
		return fmt.Errorf("delete refresh session: %w", err)
	}

	return nil
}
