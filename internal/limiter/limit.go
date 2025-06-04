package limiter

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter    *rate.Limiter
	blockedAt  time.Time
	blockUntil time.Time
}

type Limiter struct {
	ipLimit       int
	tokenLimit    int
	blockDuration time.Duration
	ipStore       map[string]*limiterEntry
	tokenStore    map[string]*limiterEntry
	useRedis      bool
	redisClient   *redis.Client
	mu            sync.Mutex
}

func NewLimiter(ipLimit, tokenLimit, blockSeconds int) *Limiter {
	return &Limiter{
		ipLimit:       ipLimit,
		tokenLimit:    tokenLimit,
		blockDuration: time.Duration(blockSeconds) * time.Second,
		ipStore:       make(map[string]*limiterEntry),
		tokenStore:    make(map[string]*limiterEntry),
		useRedis:      false,
	}
}

func (l *Limiter) UseRedis(redisAddr, redisPassword string) {
	l.redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	l.useRedis = true
}

func (l *Limiter) Allow(ip, token string) bool {
	key := ""
	limit := 0

	if token != "" {
		key = "token:" + token
		limit = l.tokenLimit
	} else {
		key = "ip:" + ip
		limit = l.ipLimit
	}

	if l.useRedis {
		return l.allowWithRedis(key, limit)
	}

	if token != "" {
		return l.allowWithStore(l.tokenStore, key, limit)
	}
	return l.allowWithStore(l.ipStore, key, limit)
}

func (l *Limiter) allowWithStore(store map[string]*limiterEntry, key string, limit int) bool {
	entry := l.getEntry(store, key, limit)

	now := time.Now()
	if entry.blockUntil.After(now) {
		return false
	}

	if entry.limiter.Allow() {
		return true
	}

	entry.blockedAt = now
	entry.blockUntil = now.Add(l.blockDuration)
	return false
}

func (l *Limiter) getEntry(store map[string]*limiterEntry, key string, limit int) *limiterEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, exists := store[key]
	if !exists {
		entry = &limiterEntry{
			limiter: rate.NewLimiter(rate.Limit(limit), limit),
		}
		store[key] = entry
	}
	return entry
}

func (l *Limiter) allowWithRedis(key string, limit int) bool {
	ctx := context.Background()
	blockKey := "block:" + key
	countKey := "count:" + key + ":" + time.Now().Format("20060102150405")

	blocked, err := l.redisClient.Exists(ctx, blockKey).Result()
	if err != nil {
		log.Printf("[Redis] ERRO ao verificar %s: %v", key, err)
		return false
	}
	if blocked > 0 {
		return false
	}

	count, err := l.redisClient.Incr(ctx, countKey).Result()
	if err != nil {
		log.Printf("[Redis] ERRO ao verificar %s: %v", key, err)
		return false
	}
	if count == 1 {
		_ = l.redisClient.Expire(ctx, countKey, time.Second)
	}
	if count > int64(limit) {
		_ = l.redisClient.Set(ctx, blockKey, "1", l.blockDuration).Err()
		return false
	}
	return true
}
