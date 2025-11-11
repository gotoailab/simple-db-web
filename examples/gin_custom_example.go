//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/chenhg5/simple-db-web/handlers"
	"github.com/gin-gonic/gin"
)

// 这是一个使用自定义 Gin 引擎的示例
// 可以添加中间件、自定义配置等
// 运行方式: go run examples/gin_custom_example.go
func main() {
	// 创建服务器实例
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 创建自定义 Gin 引擎
	gin.SetMode(gin.ReleaseMode) // 设置为发布模式
	engine := gin.New()

	// 添加中间件（可选）
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 可以添加自定义中间件
	engine.Use(func(c *gin.Context) {
		// 添加 CORS 头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Connection-ID")
		c.Next()
	})

	// 创建 Gin 适配器
	ginRouter := handlers.NewGinRouter(engine)

	// 注册路由到 Gin
	server.RegisterRoutes(ginRouter)

	// 启动服务器
	addr := ":8080"
	log.Printf("Gin 服务器启动在 %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
