package middleware

import (
	"time"

	"md2html/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AccessLog 请求访问日志中间件
// 记录每个请求的方法、路径、状态码、耗时、客户端IP
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		if query != "" {
			path = path + "?" + query
		}

		if status >= 500 {
			logger.Error("[HTTP] %s %s %d %v %s", method, path, status, latency, clientIP)
		} else if status >= 400 {
			logger.Warn("[HTTP] %s %s %d %v %s", method, path, status, latency, clientIP)
		} else {
			logger.Info("[HTTP] %s %s %d %v %s", method, path, status, latency, clientIP)
		}
	}
}
