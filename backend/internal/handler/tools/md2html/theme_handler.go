package md2html

import (
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/pkg/response"

	"github.com/gin-gonic/gin"
)

type ThemeHandler struct {
	themeService md2htmlService.ThemeService
}

func NewThemeHandler(themeService md2htmlService.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}

func (h *ThemeHandler) GetThemes(c *gin.Context) {
	themes := h.themeService.GetThemes()
	response.Success(c, themes)
}
