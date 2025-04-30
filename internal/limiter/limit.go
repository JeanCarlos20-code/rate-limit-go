package limiter

import (
	"sync"
	"time"

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
	mu            sync.Mutex
}

func NewLimiter(ipLimit, tokenLimit, blockSeconds int) *Limiter {
	return &Limiter{
		ipLimit:       ipLimit,
		tokenLimit:    tokenLimit,
		blockDuration: time.Duration(blockSeconds) * time.Second,
		ipStore:       make(map[string]*limiterEntry),
		tokenStore:    make(map[string]*limiterEntry),
	}
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

func (l *Limiter) Allow(ip, token string) bool {
	if token != "" {
		return l.allowWithStore(l.tokenStore, "token:"+token, l.tokenLimit)
	}
	return l.allowWithStore(l.ipStore, "ip:"+ip, l.ipLimit)
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
