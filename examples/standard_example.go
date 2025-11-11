//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/chenhg5/simple-db-web/handlers"
)

// 这是使用标准库 net/http 的示例（原有方式，保持兼容）
// 运行方式: go run examples/standard_example.go
func main() {
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 方式1：使用原有的 SetupRoutes 方法（向后兼容）
	server.SetupRoutes()

	// 方式2：使用新的 RegisterRoutes 方法（推荐）
	// router := handlers.NewStandardRouter()
	// server.RegisterRoutes(router)

	// 启动服务器
	addr := ":8080"
	if err := server.Start(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
