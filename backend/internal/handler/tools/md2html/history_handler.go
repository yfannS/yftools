package md2html

import (
	"database/sql"
	"strconv"

	"md2html/internal/middleware"
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/internal/model"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	historyService md2htmlService.HistoryService
}

func NewHistoryHandler(historyService md2htmlService.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: historyService}
}

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
		response.InternalError(c, "获取历史记录失败")
		return
	}

	response.Success(c, model.HistoryPageResponse{
		Data:  result.Data,
		Total: result.Total,
		Page:  result.Page,
		Size:  result.Size,
	})
}

func (h *HistoryHandler) SaveHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.SaveHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	theme := req.Theme
	if theme == "" {
		theme = "default"
	}

	id, err := h.historyService.Save(userID, req.Markdown, req.HTML, theme)
	if err != nil {
		response.InternalError(c, "保存失败")
		return
	}

	response.SuccessWithMessage(c, "保存成功", gin.H{"id": id})
}

func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的记录 ID")
		return
	}

	if err := h.historyService.DeleteByIDAndUserID(id, userID); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "记录不存在")
			return
		}
		response.InternalError(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
