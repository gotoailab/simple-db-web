//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/chenhg5/simple-db-web/handlers"
)

// 这是一个使用 Gin 框架的示例
// 运行前需要安装 Gin: go get -u github.com/gin-gonic/gin
// 运行方式: go run examples/gin_example.go
func main() {
	// 创建服务器实例
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 创建 Gin 适配器
	ginRouter := handlers.NewGinRouter(nil) // nil 表示使用 gin.Default()

	// 注册路由到 Gin
	server.RegisterRoutes(ginRouter)

	// 启动服务器
	addr := ":8080"
	log.Printf("Gin 服务器启动在 %s", addr)
	if err := ginRouter.Engine().Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
