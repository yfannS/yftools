package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, statusCode int, code string, message string) {
	ErrorWithData(c, statusCode, code, message, nil)
}

func ErrorWithData(c *gin.Context, statusCode int, code string, message string, data interface{}) {
	c.Set("error_code", code)
	c.Set("error_message", message)

	payload := gin.H{
		"success": false,
		"code":    code,
		"error":   message,
	}
	if data != nil {
		payload["data"] = data
	}
	if requestID := getRequestID(c); requestID != "" {
		payload["request_id"] = requestID
	}

	c.JSON(statusCode, payload)
}

func getRequestID(c *gin.Context) string {
	requestID, _ := c.Get("request_id")
	if value, ok := requestID.(string); ok {
		return value
	}
	return ""
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeInvalidParams, message)
}

func BadRequestCode(c *gin.Context, code string, message string) {
	Error(c, http.StatusBadRequest, code, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeTokenInvalid, message)
}

func UnauthorizedCode(c *gin.Context, code string, message string) {
	Error(c, http.StatusUnauthorized, code, message)
}

func UnauthorizedCodeData(c *gin.Context, code string, message string, data interface{}) {
	ErrorWithData(c, http.StatusUnauthorized, code, message, data)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, CodeRateLimited, message)
}

func TooManyRequestsCode(c *gin.Context, code string, message string) {
	Error(c, http.StatusTooManyRequests, code, message)
}

func TooManyRequestsCodeData(c *gin.Context, code string, message string, data interface{}) {
	ErrorWithData(c, http.StatusTooManyRequests, code, message, data)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeInternalError, message)
}

func ForbiddenCode(c *gin.Context, code string, message string) {
	Error(c, http.StatusForbidden, code, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeInternalError, message)
}

func NotFoundCode(c *gin.Context, code string, message string) {
	Error(c, http.StatusNotFound, code, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, message)
}

func InternalErrorCode(c *gin.Context, code string, message string) {
	Error(c, http.StatusInternalServerError, code, message)
}
