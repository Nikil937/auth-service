package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshRepository interface {
	Save(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error
	Get(ctx context.Context, token string) (uuid.UUID, error)
	Delete(ctx context.Context, token string) error
}
