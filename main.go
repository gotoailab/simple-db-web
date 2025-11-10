package main

import (
	"dbweb/handlers"
	"log"
)

func main() {
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	server.SetupRoutes()

	// 默认监听8080端口
	addr := ":8080"
	if err := server.Start(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
