package handlers

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

// EchoRouter Echo 框架的路由适配器
type EchoRouter struct {
	echo *echo.Echo
}

// NewEchoRouter 创建 Echo 路由适配器
// 如果 e 为 nil，会创建一个新的 echo.New() 实例
func NewEchoRouter(e *echo.Echo) *EchoRouter {
	if e == nil {
		e = echo.New()
	}
	return &EchoRouter{
		echo: e,
	}
}

// GET 注册 GET 路由
func (r *EchoRouter) GET(path string, handler http.HandlerFunc) {
	r.echo.GET(path, echo.WrapHandler(http.HandlerFunc(handler)))
}

// POST 注册 POST 路由
func (r *EchoRouter) POST(path string, handler http.HandlerFunc) {
	r.echo.POST(path, echo.WrapHandler(http.HandlerFunc(handler)))
}

// Static 注册静态文件路由
func (r *EchoRouter) Static(path, dir string) {
	r.echo.Static(path, dir)
}

// StaticFS 注册静态文件路由（使用 embed.FS）
func (r *EchoRouter) StaticFS(path string, fsys fs.FS) {
	// staticFS 使用 all:static，所以路径是 static/
	subFS, err := fs.Sub(fsys, "static")
	if err != nil {
		// 如果失败，尝试直接使用
		subFS = fsys
	}
	r.echo.StaticFS(path, subFS)
}

// HandleFunc 注册任意 HTTP 方法的路由
func (r *EchoRouter) HandleFunc(path string, handler http.HandlerFunc) {
	r.echo.Any(path, echo.WrapHandler(http.HandlerFunc(handler)))
}

// Echo 返回 Echo 实例，用于启动服务器
func (r *EchoRouter) Echo() *echo.Echo {
	return r.echo
}

// SetPrefix 设置路由前缀（Echo 使用路由组实现）
func (r *EchoRouter) SetPrefix(prefix string) {
	// Echo 通过路由组实现前缀，这里不做处理
	// 使用 NewPrefixRouter 包装器来添加前缀支持
}

// GetPrefix 获取当前路由前缀
func (r *EchoRouter) GetPrefix() string {
	return ""
}
