package router

import (
	"md2html/internal/config"
	"md2html/internal/handler"
	md2htmlHandler "md2html/internal/handler/tools/md2html"
	"md2html/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	rateLimitCfg config.RateLimitConfig,
	authHandler *handler.AuthHandler,
	convertHandler *md2htmlHandler.ConvertHandler,
	historyHandler *md2htmlHandler.HistoryHandler,
	themeHandler *md2htmlHandler.ThemeHandler,
) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.AccessLog())
	r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/api/health", handler.HealthCheck)

	// ====== 认证路由（跨工具共享） ======
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", middleware.RegisterRateLimit(rateLimitCfg), authHandler.Register)
		auth.POST("/login", middleware.LoginRateLimit(rateLimitCfg), authHandler.Login)
		auth.GET("/profile", middleware.Auth(), authHandler.GetProfile)
		auth.POST("/logout", middleware.Auth(), authHandler.Logout)
	}

	// ====== 工具路由（按工具分组） ======
	toolsGroup := r.Group("/api/tools")

	// --- md2html 工具 ---
	md2html := toolsGroup.Group("/md2html")
	{
		// 公开接口
		md2html.POST("/convert", convertHandler.ConvertMarkdown)

		// 需认证接口
		md2html.GET("/history", middleware.Auth(), historyHandler.GetHistory)           // 列表（轻量）
		md2html.GET("/history/:id", middleware.Auth(), historyHandler.GetHistoryDetail) // 详情
		md2html.POST("/history", middleware.Auth(), historyHandler.SaveHistory)         // 保存
		md2html.PATCH("/history/:id", middleware.Auth(), historyHandler.RenameHistory) // 改标题
		md2html.DELETE("/history/:id", middleware.Auth(), historyHandler.DeleteHistory) // 删除

		// 主题（公开）
		md2html.GET("/themes", themeHandler.GetThemes)
	}

	// 未来工具在此扩展：
	// jsonTools := toolsGroup.Group("/json-formatter")
	// {
	//     jsonTools.POST("/format", ...)
	// }

	return r
}
