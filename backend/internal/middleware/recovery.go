package middleware

import (
	"net/http"

	"md2html/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("[PANIC] %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}
