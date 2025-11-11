package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GinRouter Gin 框架的路由适配器
type GinRouter struct {
	engine *gin.Engine
}

// NewGinRouter 创建 Gin 路由适配器
// 如果 engine 为 nil，会创建一个新的 gin.Default() 实例
func NewGinRouter(engine *gin.Engine) *GinRouter {
	if engine == nil {
		engine = gin.Default()
	}
	return &GinRouter{
		engine: engine,
	}
}

// GET 注册 GET 路由
func (r *GinRouter) GET(path string, handler http.HandlerFunc) {
	r.engine.GET(path, gin.WrapH(http.HandlerFunc(handler)))
}

// POST 注册 POST 路由
func (r *GinRouter) POST(path string, handler http.HandlerFunc) {
	r.engine.POST(path, gin.WrapH(http.HandlerFunc(handler)))
}

// Static 注册静态文件路由
func (r *GinRouter) Static(path, dir string) {
	r.engine.Static(path, dir)
}

// HandleFunc 注册任意 HTTP 方法的路由（Gin 中会注册为 GET 和 POST）
func (r *GinRouter) HandleFunc(path string, handler http.HandlerFunc) {
	r.engine.Any(path, gin.WrapH(http.HandlerFunc(handler)))
}

// Engine 返回 Gin 引擎实例，用于启动服务器
func (r *GinRouter) Engine() *gin.Engine {
	return r.engine
}

// SetPrefix 设置路由前缀（Gin 使用路由组实现）
func (r *GinRouter) SetPrefix(prefix string) {
	// Gin 通过路由组实现前缀，这里不做处理
	// 使用 NewPrefixRouter 包装器来添加前缀支持
}

// GetPrefix 获取当前路由前缀
func (r *GinRouter) GetPrefix() string {
	return ""
}

