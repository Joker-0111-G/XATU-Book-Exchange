package main

import (
	"fmt"
	"log"

	"xatu-book-exchange/config"
	"xatu-book-exchange/database"
	"xatu-book-exchange/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 2. 连接 MySQL
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}

	// 3. 连接 Redis（非必须，失败仅警告）
	if err := database.InitRedis(); err != nil {
		log.Printf("⚠️  Redis 连接失败（不影响启动）: %v", err)
	}

	// 4. 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 5. 初始化路由
	r := routes.SetupRouter()

	// 6. 启动服务
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("🚀  XATU Book Exchange 启动成功，监听 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
