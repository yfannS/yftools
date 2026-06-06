package middleware

import (
	"strings"

	appJwt "md2html/pkg/jwt"
	"md2html/pkg/logger"
	appRedis "md2html/pkg/redis"
	"md2html/pkg/response"
	"md2html/pkg/session"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 1. 先验证 JWT 签名和有效期
		claims, err := appJwt.ParseToken(tokenString)
		if err != nil {
			logger.Warn("[Auth] JWT parse failed: %v", err)
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		// 2. 如果 Redis 可用，校验会话是否存在（支持主动登出/踢下线）
		if appRedis.IsAvailable() {
			sess, err := session.Get(tokenString)
			if err != nil {
				logger.Warn("[Auth] Redis session check error: %v", err)
				// Redis 读取异常，降级为 JWT 继续放行
			} else if sess == nil {
				// 会话不存在 = 已被登出或过期
				logger.Warn("[Auth] Session not found in Redis for userID=%d", claims.UserID)
				response.Unauthorized(c, "登录已失效，请重新登录")
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	id, _ := c.Get("user_id")
	if v, ok := id.(int64); ok {
		return v
	}
	return 0
}

func GetUsername(c *gin.Context) string {
	name, _ := c.Get("username")
	if v, ok := name.(string); ok {
		return v
	}
	return ""
}
