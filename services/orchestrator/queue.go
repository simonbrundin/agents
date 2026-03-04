package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func pushJob(ctx context.Context, payload string) error {
	r := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	defer r.Close()

	stream := "agent-jobs"
	_, err := r.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: map[string]interface{}{"payload": payload}}).Result()
	if err != nil {
		return fmt.Errorf("xadd failed: %w", err)
	}
	return nil
}
