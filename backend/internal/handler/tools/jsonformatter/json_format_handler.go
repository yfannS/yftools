package jsonformatter

import (
	"database/sql"
	"strconv"

	"md2html/internal/middleware"
	"md2html/internal/model"
	jsonService "md2html/internal/service/tools/jsonformatter"
	"md2html/pkg/logger"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type JsonFormatHandler struct {
	jsonService jsonService.JsonFormatService
}

func NewJsonFormatHandler(jsonService jsonService.JsonFormatService) *JsonFormatHandler {
	return &JsonFormatHandler{jsonService: jsonService}
}

// Format 格式化/压缩 JSON
func (h *JsonFormatHandler) Format(c *gin.Context) {
	var req model.JsonFormatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	indent := req.Indent
	if indent <= 0 {
		indent = 2
	}

	result, err := h.jsonService.Format(req.Input, indent, req.Minify)
	if err != nil {
		logger.WarnKV("[JsonFormat] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"error":      err,
		})
		response.BadRequestCode(c, response.CodeJsonFormatFailed, err.Error())
		return
	}

	logger.Debug("[JsonFormat] success size=%d minified=%v", result.Size, result.Minified)
	response.Success(c, result)
}

// Validate 校验 JSON
func (h *JsonFormatHandler) Validate(c *gin.Context) {
	var req model.JsonValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	result, err := h.jsonService.Validate(req.Input)
	if err != nil {
		logger.WarnKV("[JsonValidate] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeJsonValidateFailed, "校验失败")
		return
	}

	logger.Debug("[JsonValidate] valid=%v keys=%d depth=%d", result.Valid, result.Keys, result.Depth)
	response.Success(c, result)
}

// GetHistory 获取 JSON 格式化历史列表
func (h *JsonFormatHandler) GetHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.jsonService.ListHistory(userID, page, pageSize)
	if err != nil {
		logger.ErrorKV("[JsonGetHistory] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeJsonHistoryFailed, "获取历史记录失败")
		return
	}

	response.Success(c, model.JsonHistoryPageResponse{
		Data:       result.Data,
		Total:      result.Total,
		Page:       result.Page,
		Size:       result.Size,
		TotalPages: result.TotalPages,
	})
}

// GetHistoryDetail 获取 JSON 格式化历史详情
func (h *JsonFormatHandler) GetHistoryDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequestCode(c, response.CodeInvalidRecordID, "无效的记录 ID")
		return
	}

	detail, err := h.jsonService.GetHistoryDetail(id, userID)
	if err != nil {
		logger.WarnKV("[JsonGetHistoryDetail] not found", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"history_id": id,
		})
		response.NotFoundCode(c, response.CodeHistoryNotFound, "记录不存在")
		return
	}

	response.Success(c, detail)
}

// SaveHistory 保存 JSON 格式化历史
func (h *JsonFormatHandler) SaveHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.JsonSaveHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	id, err := h.jsonService.SaveHistory(userID, req.Input, req.Output)
	if err != nil {
		logger.ErrorKV("[JsonSaveHistory] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeJsonHistoryFailed, "保存失败")
		return
	}

	logger.Debug("[JsonSaveHistory] success id=%d", id)
	response.SuccessWithMessage(c, "保存成功", gin.H{"id": id})
}

// DeleteHistory 删除 JSON 格式化历史
func (h *JsonFormatHandler) DeleteHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequestCode(c, response.CodeInvalidRecordID, "无效的记录 ID")
		return
	}

	if err := h.jsonService.DeleteHistory(id, userID); err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundCode(c, response.CodeHistoryNotFound, "记录不存在")
			return
		}
		logger.ErrorKV("[JsonDeleteHistory] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"history_id": id,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeJsonHistoryFailed, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// RenameHistory 修改 JSON 格式化历史标题
func (h *JsonFormatHandler) RenameHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequestCode(c, response.CodeInvalidRecordID, "无效的记录 ID")
		return
	}

	var req model.RenameHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	if err := h.jsonService.RenameHistory(id, userID, req.Title); err != nil {
		logger.ErrorKV("[JsonRenameHistory] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"history_id": id,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeHistoryRenameFailed, "修改标题失败")
		return
	}

	response.SuccessWithMessage(c, "修改成功", nil)
}
