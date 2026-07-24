// Package cache provides a small key/value store with TTL and atomic counters.
// It uses Redis when configured, otherwise an in-process map (dev/single-node).
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store is the minimal interface the app depends on.
type Store interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, bool, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string, ttl time.Duration) (int64, error)
}

// --- Redis-backed ---

type redisStore struct{ rdb *redis.Client }

// NewRedis builds a Redis-backed store.
func NewRedis(addr, password string, db int) (Store, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &redisStore{rdb: rdb}, nil
}

func (s *redisStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.rdb.Set(ctx, key, value, ttl).Err()
}

func (s *redisStore) Get(ctx context.Context, key string) (string, bool, error) {
	v, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return v, true, nil
}

func (s *redisStore) Del(ctx context.Context, key string) error {
	return s.rdb.Del(ctx, key).Err()
}

func (s *redisStore) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	n, err := s.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		_ = s.rdb.Expire(ctx, key, ttl).Err()
	}
	return n, nil
}

// --- In-memory fallback ---

type entry struct {
	value     string
	count     int64
	expiresAt time.Time
}

type memStore struct {
	mu sync.Mutex
	m  map[string]*entry
}

// NewMemory builds an in-process store (not for multi-node production).
func NewMemory() Store {
	s := &memStore{m: map[string]*entry{}}
	go s.reap()
	return s
}

func (s *memStore) reap() {
	t := time.NewTicker(time.Minute)
	for range t.C {
		s.mu.Lock()
		now := time.Now()
		for k, e := range s.m {
			if now.After(e.expiresAt) {
				delete(s.m, k)
			}
		}
		s.mu.Unlock()
	}
}

func (s *memStore) Set(_ context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = &entry{value: value, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (s *memStore) Get(_ context.Context, key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[key]
	if !ok || time.Now().After(e.expiresAt) {
		return "", false, nil
	}
	return e.value, true, nil
}

func (s *memStore) Del(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
	return nil
}

func (s *memStore) Incr(_ context.Context, key string, ttl time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[key]
	if !ok || time.Now().After(e.expiresAt) {
		e = &entry{expiresAt: time.Now().Add(ttl)}
		s.m[key] = e
	}
	e.count++
	return e.count, nil
}
