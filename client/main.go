package main

import (
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/chenhg5/simple-db-web/handlers"
	"github.com/gin-gonic/gin"
)

//go:embed static/user-management.js
var userManagementScript string

//go:embed templates/login.html
var loginHTML string

func main() {
	// 初始化数据库
	dbPath := "client.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	if err := InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer CloseDB()

	// 定期清理过期会话
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := CleanExpiredSessions(); err != nil {
				log.Printf("Failed to clean expired sessions: %v", err)
			}
		}
	}()

	// 创建核心服务器
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 设置自定义脚本（注入用户管理UI）
	server.SetCustomScript(userManagementScript)

	// 创建Gin引擎
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 登录页面模板已通过 embed 嵌入，不需要从文件系统加载

	// 注册认证中间件
	engine.Use(AuthMiddleware())

	// 登录相关路由（不需要认证）
	engine.GET("/login", LoginPage)
	engine.POST("/api/auth/login", LoginAPI)

	// 认证后的路由
	engine.GET("/api/auth/current", GetCurrentUserAPI)
	engine.POST("/api/auth/logout", LogoutAPI)
	engine.POST("/api/auth/password", UpdatePasswordAPI)

	// 用户管理路由（需要管理员权限）
	adminGroup := engine.Group("/api/users")
	adminGroup.Use(AdminMiddleware())
	{
		adminGroup.GET("", GetUsersAPI)
		adminGroup.POST("", CreateUserAPI)
		adminGroup.PUT("/:id", UpdateUserAPI)
		adminGroup.DELETE("/:id", DeleteUserAPI)
	}

	// 注册核心服务器的路由
	ginRouter := handlers.NewGinRouter(engine)
	server.RegisterRoutes(ginRouter)

	// 启动服务器
	port := ":8080"
	if len(os.Args) > 2 {
		port = ":" + os.Args[2]
	}

	log.Printf("Server starting on http://localhost%s", port)
	log.Printf("Default admin account: admin / admin123")
	log.Printf("Database file: %s", dbPath)

	if err := engine.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
