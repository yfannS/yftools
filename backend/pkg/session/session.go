package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	appRedis "md2html/pkg/redis"

	"github.com/redis/go-redis/v9"
)

const (
	// sessionKeyPrefix Redis key 前缀，格式：项目:模块:业务:标识
	sessionKeyPrefix = "yftools:auth:session:"
	// defaultTTL 默认过期时间 7 天（与 JWT 一致）
	defaultTTL = 7 * 24 * time.Hour
)

// Session 登录会话信息
type Session struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// Set 存储登录会话到 Redis
// key 为 token 字符串，value 为 Session JSON
func Set(token string, session *Session, ttl time.Duration) error {
	client := appRedis.GetClient()
	if client == nil {
		return fmt.Errorf("redis not available")
	}

	if ttl <= 0 {
		ttl = defaultTTL
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Set(ctx, sessionKeyPrefix+token, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set session: %w", err)
	}

	return nil
}

// Get 从 Redis 获取登录会话
// 返回 nil 表示会话不存在（已过期或已登出）
func Get(token string) (*Session, error) {
	client := appRedis.GetClient()
	if client == nil {
		return nil, fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := client.Get(ctx, sessionKeyPrefix+token).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get session: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, nil
}

// Delete 删除登录会话（登出）
func Delete(token string) error {
	client := appRedis.GetClient()
	if client == nil {
		return fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return client.Del(ctx, sessionKeyPrefix+token).Err()
}

// Exists 检查会话是否存在
func Exists(token string) (bool, error) {
	client := appRedis.GetClient()
	if client == nil {
		return false, fmt.Errorf("redis not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	n, err := client.Exists(ctx, sessionKeyPrefix+token).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
