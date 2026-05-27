package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// New initializes a Redis client for session/presence storage with startup retry logic.
func New(ctx context.Context, log *slog.Logger, addr string) (*redis.Client, error) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}
	opts.PoolSize = 20

	client := redis.NewClient(opts)

	var pingErr error
	for attempt := 1; attempt <= 5; attempt++ {
		pingErr = client.Ping(ctx).Err()
		if pingErr == nil {
			log.Info("redis connected successfully")
			return client, nil
		}
		log.Warn("redis ping failed, retrying...", slog.Int("attempt", attempt), slog.String("error", pingErr.Error()))
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}

	return nil, fmt.Errorf("failed to connect to redis after retries: %w", pingErr)
}
