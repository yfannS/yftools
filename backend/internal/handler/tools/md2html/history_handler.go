package md2html

import (
	"database/sql"
	"strconv"

	"md2html/internal/middleware"
	"md2html/internal/model"
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/pkg/logger"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	historyService md2htmlService.HistoryService
}

func NewHistoryHandler(historyService md2htmlService.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: historyService}
}

// GetHistory 获取历史记录列表（轻量，不含 markdown/html）
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.historyService.ListByUserID(userID, page, pageSize)
	if err != nil {
		logger.ErrorKV("[GetHistory] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"page":       page,
			"page_size":  pageSize,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeHistoryListFailed, "获取历史记录失败")
		return
	}

	logger.Debug("[GetHistory] userID=%d page=%d total=%d", userID, page, result.Total)
	response.Success(c, model.HistoryPageResponse{
		Data:       result.Data,
		Total:      result.Total,
		Page:       result.Page,
		Size:       result.Size,
		TotalPages: result.TotalPages,
	})
}

// GetHistoryDetail 获取历史记录详情（含完整 markdown）
func (h *HistoryHandler) GetHistoryDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequestCode(c, response.CodeInvalidRecordID, "无效的记录 ID")
		return
	}

	detail, err := h.historyService.GetDetail(id, userID)
	if err != nil {
		logger.WarnKV("[GetHistoryDetail] not found", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"history_id": id,
		})
		response.NotFoundCode(c, response.CodeHistoryNotFound, "记录不存在")
		return
	}

	logger.Debug("[GetHistoryDetail] id=%d title=%s", id, detail.Title)
	response.Success(c, detail)
}

// SaveHistory 保存历史记录（html 可选，title 由后端自动从 markdown 提取）
func (h *HistoryHandler) SaveHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.SaveHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	logger.InfoKV("[SaveHistory] received", logger.Fields{
		"request_id":   middleware.GetRequestID(c),
		"user_id":      userID,
		"markdown_len": len(req.Markdown),
		"html_len":     len(req.HTML),
		"theme":        req.Theme,
	})

	theme := req.Theme
	if theme == "" {
		theme = "default"
	}

	id, err := h.historyService.Save(userID, req.Markdown, req.HTML, theme)
	if err != nil {
		logger.ErrorKV("[SaveHistory] save failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"theme":      theme,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeHistorySaveFailed, "保存失败")
		return
	}

	logger.Debug("[SaveHistory] save success, id=%d", id)
	response.SuccessWithMessage(c, "保存成功", gin.H{"id": id})
}

// DeleteHistory 软删除历史记录
func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequestCode(c, response.CodeInvalidRecordID, "无效的记录 ID")
		return
	}

	if err := h.historyService.DeleteByIDAndUserID(id, userID); err != nil {
		if err == sql.ErrNoRows {
			logger.WarnKV("[DeleteHistory] not found", logger.Fields{
				"request_id": middleware.GetRequestID(c),
				"user_id":    userID,
				"history_id": id,
			})
			response.NotFoundCode(c, response.CodeHistoryNotFound, "记录不存在")
			return
		}
		logger.ErrorKV("[DeleteHistory] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"user_id":    userID,
			"history_id": id,
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeHistoryDeleteFailed, "删除失败")
		return
	}

	logger.InfoKV("[DeleteHistory] success", logger.Fields{
		"request_id": middleware.GetRequestID(c),
		"user_id":    userID,
		"history_id": id,
	})
	response.SuccessWithMessage(c, "删除成功", nil)
}
