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
			response.UnauthorizedCode(c, response.CodeTokenRequired, "未提供认证令牌")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.UnauthorizedCode(c, response.CodeTokenMalformed, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 1. 先验证 JWT 签名和有效期
		claims, err := appJwt.ParseToken(tokenString)
		if err != nil {
			logger.WarnKV("[Auth] JWT parse failed", logger.Fields{
				"request_id": GetRequestID(c),
				"client_ip":  c.ClientIP(),
				"error":      err,
			})
			response.UnauthorizedCode(c, response.CodeTokenInvalid, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		// 2. 如果 Redis 可用，校验会话是否存在（支持主动登出/踢下线）
		if appRedis.IsAvailable() {
			sess, err := session.Get(tokenString)
			if err != nil {
				logger.WarnKV("[Auth] Redis session check error", logger.Fields{
					"request_id": GetRequestID(c),
					"user_id":    claims.UserID,
					"error":      err,
				})
				// Redis 读取异常，降级为 JWT 继续放行
			} else if sess == nil {
				// 会话不存在 = 已被登出或过期
				logger.WarnKV("[Auth] Session not found", logger.Fields{
					"request_id": GetRequestID(c),
					"user_id":    claims.UserID,
				})
				response.UnauthorizedCode(c, response.CodeSessionExpired, "登录已失效，请重新登录")
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
