//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/gotoailab/simple-db-web/handlers"
)

// 这是一个使用 Echo 框架的示例
// 运行前需要安装 Echo: go get -u github.com/labstack/echo/v4
// 运行方式: go run examples/echo_example.go
func main() {
	// 创建服务器实例
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 创建 Echo 适配器
	echoRouter := handlers.NewEchoRouter(nil) // nil 表示使用 echo.New()

	// 注册路由到 Echo
	server.RegisterRoutes(echoRouter)

	// 启动服务器
	addr := ":8080"
	log.Printf("Echo 服务器启动在 %s", addr)
	if err := echoRouter.Echo().Start(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
