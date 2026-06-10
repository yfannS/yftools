package main

import (
	"md2html/internal/config"
	"md2html/internal/handler"
	md2htmlHandler "md2html/internal/handler/tools/md2html"
	"md2html/internal/repository"
	md2htmlRepo "md2html/internal/repository/tools/md2html"
	"md2html/internal/router"
	"md2html/internal/service"
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/pkg/converter"
	appJwt "md2html/pkg/jwt"
	"md2html/pkg/logger"
	appRedis "md2html/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 初始化日志
	logger.SetLevel(logger.ParseLevel(cfg.Server.LogLevel))
	if err := logger.AddFileOutput("logs/app.log"); err != nil {
		logger.Warn("Failed to add log file output: %v", err)
	}

	// 初始化 JWT 配置
	appJwt.SetSecret(cfg.JWT.Secret)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db := cfg.InitDB()
	defer db.Close()

	// 初始化 Redis
	cfg.InitRedis()
	defer appRedis.Close()

	// ====== Repository 层 ======
	userRepo := repository.NewUserRepository(db)
	historyRepo := md2htmlRepo.NewHistoryRepository(db)

	// ====== Service 层 ======
	userService := service.NewUserService(userRepo, cfg.JWT.Expire, cfg.RateLimit)
	goldmarkConverter := converter.NewGoldmarkConverter()
	convertService := md2htmlService.NewConvertService(goldmarkConverter)
	historyService := md2htmlService.NewHistoryService(historyRepo)
	themeService := md2htmlService.NewThemeService()

	// ====== Handler 层 ======
	authHandler := handler.NewAuthHandler(userService)
	convertHandler := md2htmlHandler.NewConvertHandler(convertService)
	historyHandler := md2htmlHandler.NewHistoryHandler(historyService)
	themeHandler := md2htmlHandler.NewThemeHandler(themeService)

	// ====== 路由 ======
	r := router.SetupRouter(cfg.RateLimit, authHandler, convertHandler, historyHandler, themeHandler)

	logger.Info("Server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}
