// limiter/redis_store.go
package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		}),
	}
}

func (r *RedisStore) Allow(ctx context.Context, key string, limit int, blockSeconds int) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	countKey := fmt.Sprintf("count:%s:%d", key, time.Now().Unix())

	blocked, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}
	if blocked > 0 {
		return false, nil
	}

	count, err := r.client.Incr(ctx, countKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		_ = r.client.Expire(ctx, countKey, time.Second).Err()
	}

	if count > int64(limit) {
		_ = r.client.Set(ctx, blockKey, "1", time.Duration(blockSeconds)*time.Second).Err()
		return false, nil
	}

	return true, nil
}
