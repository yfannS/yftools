package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"md2html/internal/config"
	"md2html/pkg/logger"
	"md2html/pkg/ratelimit"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type authRequestMeta struct {
	IP       string
	Username string
}

func RegisterRateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	rules := []limitRule{
		{
			name:   "register_ip_per_minute",
			config: cfg.Register.IPPerMinute,
			key: func(meta authRequestMeta, prefix string) string {
				if meta.IP == "" {
					return ""
				}
				return fmt.Sprintf("%s:register:ip_minute:%s", prefix, meta.IP)
			},
		},
		{
			name:   "register_ip_per_hour",
			config: cfg.Register.IPPerHour,
			key: func(meta authRequestMeta, prefix string) string {
				if meta.IP == "" {
					return ""
				}
				return fmt.Sprintf("%s:register:ip_hour:%s", prefix, meta.IP)
			},
		},
		{
			name:   "register_username_per_hour",
			config: cfg.Register.UsernamePerHour,
			key: func(meta authRequestMeta, prefix string) string {
				if meta.Username == "" {
					return ""
				}
				return fmt.Sprintf("%s:register:username_hour:%s", prefix, normalizeKeyPart(meta.Username))
			},
		},
	}

	return authRateLimit("register", cfg, rules)
}

func LoginRateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	rules := []limitRule{
		{
			name:   "login_ip_per_minute",
			config: cfg.Login.IPPerMinute,
			key: func(meta authRequestMeta, prefix string) string {
				if meta.IP == "" {
					return ""
				}
				return fmt.Sprintf("%s:login:ip_minute:%s", prefix, meta.IP)
			},
		},
		{
			name:   "login_username_per_15m",
			config: cfg.Login.UsernamePer15m,
			key: func(meta authRequestMeta, prefix string) string {
				if meta.Username == "" {
					return ""
				}
				return fmt.Sprintf("%s:login:username_15m:%s", prefix, normalizeKeyPart(meta.Username))
			},
		},
	}

	return authRateLimit("login", cfg, rules)
}

type limitRule struct {
	name   string
	config config.LimitRuleConfig
	key    func(meta authRequestMeta, prefix string) string
}

func authRateLimit(action string, cfg config.RateLimitConfig, rules []limitRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}

		meta := extractAuthRequestMeta(c)
		prefix := cfg.KeyPrefix()

		for _, rule := range rules {
			if !rule.config.IsEnabled() {
				continue
			}

			key := rule.key(meta, prefix)
			if key == "" {
				continue
			}

			result, err := ratelimit.Allow(key, rule.config.Limit, rule.config.Duration())
			if err != nil {
				logger.Warn("[RateLimit] %s skipped: %v", rule.name, err)
				continue
			}

			if result.Allowed {
				continue
			}

			retryAfter := secondsFromTTL(result.TTL)
			if retryAfter > 0 {
				c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			}

			logger.Warn("[RateLimit] blocked action=%s rule=%s ip=%s username=%s count=%d limit=%d",
				action, rule.name, meta.IP, meta.Username, result.Count, rule.config.Limit)
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractAuthRequestMeta(c *gin.Context) authRequestMeta {
	meta := authRequestMeta{
		IP: c.ClientIP(),
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return meta
	}
	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	if len(bodyBytes) == 0 {
		return meta
	}

	var payload struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(bodyBytes, &payload); err == nil {
		meta.Username = strings.TrimSpace(payload.Username)
	}

	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	return meta
}

func normalizeKeyPart(value string) string {
	return url.QueryEscape(strings.ToLower(strings.TrimSpace(value)))
}

func secondsFromTTL(ttl time.Duration) int64 {
	if ttl <= 0 {
		return 0
	}
	seconds := int64(ttl / time.Second)
	if ttl%time.Second != 0 {
		seconds++
	}
	return seconds
}
