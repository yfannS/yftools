package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"md2html/pkg/logger"
)

var client *redis.Client

// Config Redis 配置
type Config struct {
	Host string
	Port int
}

// Init 初始化 Redis 连接
func Init(cfg Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	logger.Info("Redis connected: %s:%d", cfg.Host, cfg.Port)
	return nil
}

// GetClient 获取 Redis 客户端（供 session store 等使用）
func GetClient() *redis.Client {
	return client
}

// Close 关闭 Redis 连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// IsAvailable 检查 Redis 是否可用
func IsAvailable() bool {
	if client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Ping(ctx).Err() == nil
}
