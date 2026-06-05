package md2html

import (
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/internal/model"
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
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	rawHTML, formattedHTML, err := h.convertService.Convert(req.Markdown)
	if err != nil {
		response.InternalError(c, "转换失败: "+err.Error())
		return
	}

	response.Success(c, model.ConvertResponse{
		HTML:      rawHTML,
		Formatted: formattedHTML,
	})
}
