package md2html

import "md2html/internal/model"

type ThemeService interface {
	GetThemes() []model.Theme
}

type themeService struct{}

func NewThemeService() ThemeService {
	return &themeService{}
}

func (s *themeService) GetThemes() []model.Theme {
	return []model.Theme{
		{ID: "default", Name: "默认主题", Description: "简洁优雅的默认样式"},
		{ID: "github", Name: "GitHub 风格", Description: "GitHub Markdown 渲染风格"},
		{ID: "dark", Name: "暗色主题", Description: "深色背景，护眼风格"},
	}
}
