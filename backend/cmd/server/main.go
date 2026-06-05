package main

import (
	"log"

	"md2html/internal/config"
	"md2html/internal/handler"
	md2htmlHandler "md2html/internal/handler/tools/md2html"
	"md2html/internal/repository"
	md2htmlRepo "md2html/internal/repository/tools/md2html"
	"md2html/internal/router"
	"md2html/internal/service"
	md2htmlService "md2html/internal/service/tools/md2html"
	"md2html/pkg/converter"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db := cfg.InitDB()
	defer db.Close()

	// 自动迁移
	repository.RunMigrations(db)

	// ====== Repository 层 ======
	userRepo := repository.NewUserRepository(db)
	historyRepo := md2htmlRepo.NewHistoryRepository(db)

	// ====== Service 层 ======
	userService := service.NewUserService(userRepo, cfg.JWT.Expire)
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
	r := router.SetupRouter(authHandler, convertHandler, historyHandler, themeHandler)

	log.Printf("Server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
