package middleware

import (
	"strings"

	appJwt "md2html/pkg/jwt"
	"md2html/pkg/response"

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

		claims, err := appJwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
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
