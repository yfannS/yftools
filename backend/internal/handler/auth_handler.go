package handler

import (
	"errors"
	"strings"

	"md2html/internal/middleware"
	"md2html/internal/model"
	"md2html/internal/service"
	"md2html/pkg/logger"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnKV("[Register] bad request", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"error":      err,
		})
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.Register(req); err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			logger.WarnKV("[Register] username exists", logger.Fields{
				"request_id": middleware.GetRequestID(c),
				"client_ip":  c.ClientIP(),
				"username":   req.Username,
			})
			response.BadRequestCode(c, response.CodeUsernameExists, err.Error())
			return
		}

		logger.ErrorKV("[Register] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"username":   req.Username,
			"error":      err,
		})
		response.InternalError(c, "注册失败")
		return
	}

	logger.InfoKV("[Register] success", logger.Fields{
		"request_id": middleware.GetRequestID(c),
		"client_ip":  c.ClientIP(),
		"username":   req.Username,
	})
	response.SuccessWithMessage(c, "注册成功", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnKV("[Login] bad request", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"error":      err,
		})
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	token, err := h.userService.Login(req, c.ClientIP())
	if err != nil {
		if errors.Is(err, service.ErrTooManyRequests) {
			logger.WarnKV("[Login] rate limited", logger.Fields{
				"request_id": middleware.GetRequestID(c),
				"client_ip":  c.ClientIP(),
				"username":   req.Username,
			})
			response.TooManyRequestsCode(c, response.CodeRateLimited, err.Error())
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			logger.WarnKV("[Login] invalid credentials", logger.Fields{
				"request_id": middleware.GetRequestID(c),
				"client_ip":  c.ClientIP(),
				"username":   req.Username,
			})
			response.UnauthorizedCode(c, response.CodeInvalidCredentials, err.Error())
			return
		}
		logger.ErrorKV("[Login] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"username":   req.Username,
			"error":      err,
		})
		response.InternalError(c, "登录失败")
		return
	}

	logger.InfoKV("[Login] success", logger.Fields{
		"request_id": middleware.GetRequestID(c),
		"client_ip":  c.ClientIP(),
		"username":   req.Username,
	})
	response.Success(c, model.LoginResponse{
		Token:    token,
		Username: req.Username,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		logger.WarnKV("[GetProfile] user not found", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"error":      err,
		})
		response.NotFoundCode(c, response.CodeUserNotFound, err.Error())
		return
	}

	response.Success(c, profile)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := ""
	if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 {
		tokenString = parts[1]
	}

	if err := h.userService.Logout(tokenString); err != nil {
		logger.WarnKV("[Logout] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"error":      err,
		})
	}

	logger.InfoKV("[Logout] success", logger.Fields{
		"request_id": middleware.GetRequestID(c),
	})
	response.SuccessWithMessage(c, "登出成功", nil)
}
