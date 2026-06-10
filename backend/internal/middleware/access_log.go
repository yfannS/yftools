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
		requestID := GetRequestID(c)
		errorCode, _ := c.Get("error_code")
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		if query != "" {
			path = path + "?" + query
		}

		fields := logger.Fields{
			"request_id": requestID,
			"method":     method,
			"path":       path,
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  clientIP,
		}
		if errorCode != nil {
			fields["error_code"] = errorCode
		}
		if userID != nil {
			fields["user_id"] = userID
		}
		if username != nil {
			fields["username"] = username
		}

		if status >= 500 {
			logger.ErrorKV("[HTTP] request completed", fields)
		} else if status >= 400 {
			logger.WarnKV("[HTTP] request completed", fields)
		} else {
			logger.InfoKV("[HTTP] request completed", fields)
		}
	}
}
