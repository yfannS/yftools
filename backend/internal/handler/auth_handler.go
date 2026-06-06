package handler

import (
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
		logger.Warn("[Register] bad request: %v", err)
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.Register(req); err != nil {
		logger.Warn("[Register] failed: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	logger.Info("[Register] user=%s registered", req.Username)
	response.SuccessWithMessage(c, "注册成功", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("[Login] bad request: %v", err)
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	token, err := h.userService.Login(req)
	if err != nil {
		logger.Warn("[Login] failed for user=%s: %v", req.Username, err)
		response.Unauthorized(c, err.Error())
		return
	}

	logger.Info("[Login] user=%s logged in", req.Username)
	response.Success(c, model.LoginResponse{
		Token:    token,
		Username: req.Username,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		logger.Warn("[GetProfile] failed for userID=%d: %v", userID, err)
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, profile)
}
