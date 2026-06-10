package md2html

import (
	"md2html/internal/middleware"
	"md2html/internal/model"
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/pkg/logger"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConvertHandler struct {
	convertService md2htmlService.ConvertService
}

func NewConvertHandler(convertService md2htmlService.ConvertService) *ConvertHandler {
	return &ConvertHandler{convertService: convertService}
}

func (h *ConvertHandler) ConvertMarkdown(c *gin.Context) {
	var req model.ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnKV("[Convert] bad request", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"error":      err,
		})
		response.BadRequestCode(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	rawHTML, formattedHTML, err := h.convertService.Convert(req.Markdown)
	if err != nil {
		logger.ErrorKV("[Convert] failed", logger.Fields{
			"request_id": middleware.GetRequestID(c),
			"client_ip":  c.ClientIP(),
			"error":      err,
		})
		response.InternalErrorCode(c, response.CodeConvertFailed, "转换失败: "+err.Error())
		return
	}

	logger.DebugKV("[Convert] success", logger.Fields{
		"request_id":   middleware.GetRequestID(c),
		"client_ip":    c.ClientIP(),
		"markdown_len": len(req.Markdown),
	})
	response.Success(c, model.ConvertResponse{
		HTML:      rawHTML,
		Formatted: formattedHTML,
	})
}
