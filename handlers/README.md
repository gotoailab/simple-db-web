# 路由适配器使用指南

本包提供了适配器模式，支持将 handlers 接入到不同的 Web 框架中。

## 支持的框架

- ✅ 标准库 `net/http`（默认）
- ✅ Gin
- ✅ Echo
- 🔄 其他框架（可轻松扩展）

## 快速开始

### 1. 使用标准库（原有方式，向后兼容）

```go
package main

import (
    "github.com/chenhg5/simple-db-web/handlers"
    "log"
)

func main() {
    server, err := handlers.NewServer()
    if err != nil {
        log.Fatalf("创建服务器失败: %v", err)
    }

    // 方式1：使用原有的 SetupRoutes 方法
    server.SetupRoutes()

    // 方式2：使用新的 RegisterRoutes 方法
    // router := handlers.NewStandardRouter()
    // server.RegisterRoutes(router)

    if err := server.Start(":8080"); err != nil {
        log.Fatalf("启动服务器失败: %v", err)
    }
}
```

### 2. 使用 Gin 框架

```go
package main

import (
    "github.com/chenhg5/simple-db-web/handlers"
    "log"
)

func main() {
    server, err := handlers.NewServer()
    if err != nil {
        log.Fatalf("创建服务器失败: %v", err)
    }

    // 创建 Gin 适配器
    ginRouter := handlers.NewGinRouter(nil) // nil 表示使用 gin.Default()

    // 注册路由
    server.RegisterRoutes(ginRouter)

    // 启动服务器
    if err := ginRouter.Engine().Run(":8080"); err != nil {
        log.Fatalf("启动服务器失败: %v", err)
    }
}
```

### 3. 使用自定义 Gin 引擎（添加中间件）

```go
package main

import (
    "github.com/chenhg5/simple-db-web/handlers"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    server, err := handlers.NewServer()
    if err != nil {
        log.Fatalf("创建服务器失败: %v", err)
    }

    // 创建自定义 Gin 引擎
    engine := gin.New()
    engine.Use(gin.Logger())
    engine.Use(gin.Recovery())
    
    // 添加自定义中间件
    engine.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Next()
    })

    // 创建适配器
    ginRouter := handlers.NewGinRouter(engine)

    // 注册路由
    server.RegisterRoutes(ginRouter)

    // 启动服务器
    if err := engine.Run(":8080"); err != nil {
        log.Fatalf("启动服务器失败: %v", err)
    }
}
```

### 4. 使用 Echo 框架

```go
package main

import (
    "github.com/chenhg5/simple-db-web/handlers"
    "log"
)

func main() {
    server, err := handlers.NewServer()
    if err != nil {
        log.Fatalf("创建服务器失败: %v", err)
    }

    // 创建 Echo 适配器
    echoRouter := handlers.NewEchoRouter(nil)

    // 注册路由
    server.RegisterRoutes(echoRouter)

    // 启动服务器
    if err := echoRouter.Echo().Start(":8080"); err != nil {
        log.Fatalf("启动服务器失败: %v", err)
    }
}
```

## 扩展支持其他框架

要实现其他框架的适配器，只需要实现 `Router` 接口：

```go
type Router interface {
    GET(path string, handler http.HandlerFunc)
    POST(path string, handler http.HandlerFunc)
    Static(path, dir string)
    HandleFunc(path string, handler http.HandlerFunc)
}
```

示例：为 Fiber 框架创建适配器

```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "net/http"
)

type FiberRouter struct {
    app *fiber.App
}

func NewFiberRouter(app *fiber.App) *FiberRouter {
    if app == nil {
        app = fiber.New()
    }
    return &FiberRouter{app: app}
}

func (r *FiberRouter) GET(path string, handler http.HandlerFunc) {
    r.app.Get(path, func(c *fiber.Ctx) error {
        // 将 fiber.Ctx 转换为 http.Request/ResponseWriter
        // 然后调用 handler
        return nil
    })
}

// ... 实现其他方法
```

## API 路由列表

所有 API 路由都通过 `RegisterRoutes` 方法注册：

- `GET /` - 首页
- `POST /api/connect` - 连接数据库
- `POST /api/disconnect` - 断开连接
- `GET /api/status` - 获取连接状态
- `GET /api/databases` - 获取数据库列表
- `POST /api/database/switch` - 切换数据库
- `GET /api/tables` - 获取表列表
- `GET /api/table/schema` - 获取表结构
- `GET /api/table/columns` - 获取表列信息
- `GET /api/table/data` - 获取表数据
- `POST /api/query` - 执行 SQL 查询
- `POST /api/row/update` - 更新行数据
- `POST /api/row/delete` - 删除行数据
- `GET /static/*` - 静态文件

## 注意事项

1. 所有 handler 函数都使用标准的 `http.HandlerFunc` 签名，适配器负责转换为框架特定的格式
2. 连接 ID 通过请求头 `X-Connection-ID` 传递，适配器需要确保请求头能正确传递
3. 静态文件路径默认为 `static/` 目录，可以通过修改 `RegisterRoutes` 中的路径来调整

