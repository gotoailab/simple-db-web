//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/gotoailab/simple-db-web/handlers"
)

// 这是一个使用 Gin 框架并添加路由前缀的示例
// 运行前需要安装 Gin: go get -u github.com/gin-gonic/gin
// 运行方式: go run examples/gin_with_prefix_example.go
func main() {
	// 创建服务器实例
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 创建 Gin 适配器
	ginRouter := handlers.NewGinRouter(nil)

	// 使用前缀包装器，所有路由将添加 /v1 前缀
	// 例如: /api/connect -> /v1/api/connect
	prefixedRouter := handlers.NewPrefixRouter(ginRouter, "/v1")

	// 注册路由到带前缀的适配器
	server.RegisterRoutes(prefixedRouter)

	// 启动服务器
	addr := ":8080"
	log.Printf("Gin 服务器启动在 %s，路由前缀: %s", addr, prefixedRouter.GetPrefix())
	if err := ginRouter.Engine().Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
