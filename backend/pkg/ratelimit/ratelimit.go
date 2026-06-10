package ratelimit

import (
	"context"
	"fmt"
	"time"

	appRedis "md2html/pkg/redis"

	"github.com/redis/go-redis/v9"
)

const requestTimeout = 3 * time.Second

var incrementScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
local ttl = redis.call("PTTL", KEYS[1])
local allowed = 1
if current > tonumber(ARGV[1]) then
  allowed = 0
end
return {allowed, current, ttl}
`)

type Result struct {
	Allowed bool
	Count   int64
	TTL     time.Duration
}

type CounterState struct {
	Count int64
	TTL   time.Duration
}

func Allow(key string, limit int, window time.Duration) (*Result, error) {
	if key == "" {
		return nil, fmt.Errorf("rate limit key is empty")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("rate limit must be positive")
	}
	if window <= 0 {
		return nil, fmt.Errorf("rate limit window must be positive")
	}

	client := appRedis.GetClient()
	if client == nil {
		return nil, fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	raw, err := incrementScript.Run(ctx, client, []string{key}, limit, window.Milliseconds()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis eval rate limit: %w", err)
	}

	values, ok := raw.([]interface{})
	if !ok || len(values) != 3 {
		return nil, fmt.Errorf("unexpected rate limit response: %T", raw)
	}

	return &Result{
		Allowed: toInt64(values[0]) == 1,
		Count:   toInt64(values[1]),
		TTL:     clampTTL(toInt64(values[2])),
	}, nil
}

func Increment(key string, window time.Duration) (*CounterState, error) {
	if key == "" {
		return nil, fmt.Errorf("counter key is empty")
	}
	if window <= 0 {
		return nil, fmt.Errorf("counter window must be positive")
	}

	client := appRedis.GetClient()
	if client == nil {
		return nil, fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	raw, err := incrementScript.Run(ctx, client, []string{key}, 1<<30, window.Milliseconds()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis increment counter: %w", err)
	}

	values, ok := raw.([]interface{})
	if !ok || len(values) != 3 {
		return nil, fmt.Errorf("unexpected counter response: %T", raw)
	}

	return &CounterState{
		Count: toInt64(values[1]),
		TTL:   clampTTL(toInt64(values[2])),
	}, nil
}

func Get(key string) (*CounterState, error) {
	if key == "" {
		return nil, fmt.Errorf("counter key is empty")
	}

	client := appRedis.GetClient()
	if client == nil {
		return nil, fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	count, err := client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return &CounterState{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get counter: %w", err)
	}

	ttl, err := client.PTTL(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get counter ttl: %w", err)
	}

	return &CounterState{
		Count: count,
		TTL:   normalizeTTL(ttl),
	}, nil
}

func Reset(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	client := appRedis.GetClient()
	if client == nil {
		return fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if err := client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis delete counters: %w", err)
	}

	return nil
}

func toInt64(value interface{}) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func clampTTL(ms int64) time.Duration {
	return normalizeTTL(time.Duration(ms) * time.Millisecond)
}

func normalizeTTL(ttl time.Duration) time.Duration {
	if ttl < 0 {
		return 0
	}
	return ttl
}
