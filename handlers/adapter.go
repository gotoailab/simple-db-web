package handlers

import (
	"net/http"
)

// Router 路由适配器接口，用于适配不同的 Web 框架
type Router interface {
	// GET 注册 GET 路由
	GET(path string, handler http.HandlerFunc)
	// POST 注册 POST 路由
	POST(path string, handler http.HandlerFunc)
	// Static 注册静态文件路由
	Static(path, dir string)
	// HandleFunc 注册任意 HTTP 方法的路由（用于兼容标准库）
	HandleFunc(path string, handler http.HandlerFunc)
	// SetPrefix 设置路由前缀（可选，如果适配器不支持则忽略）
	SetPrefix(prefix string)
	// GetPrefix 获取当前路由前缀
	GetPrefix() string
}

// StandardRouter 标准库 net/http 的路由适配器
type StandardRouter struct{}

// NewStandardRouter 创建标准库路由适配器
func NewStandardRouter() *StandardRouter {
	return &StandardRouter{}
}

// GET 注册 GET 路由
func (r *StandardRouter) GET(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// POST 注册 POST 路由
func (r *StandardRouter) POST(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// Static 注册静态文件路由
func (r *StandardRouter) Static(path, dir string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle(path, http.StripPrefix(path, fs))
}

// HandleFunc 注册任意 HTTP 方法的路由
func (r *StandardRouter) HandleFunc(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
}

// SetPrefix 设置路由前缀
func (r *StandardRouter) SetPrefix(prefix string) {
	// 标准库不支持前缀，使用 PrefixRouter 包装器
}

// GetPrefix 获取当前路由前缀
func (r *StandardRouter) GetPrefix() string {
	return ""
}

// PrefixRouter 路由前缀包装器，为任何 Router 添加前缀支持
type PrefixRouter struct {
	router Router
	prefix string
}

// NewPrefixRouter 创建带前缀的路由包装器
func NewPrefixRouter(router Router, prefix string) *PrefixRouter {
	// 确保前缀以 / 开头，不以 / 结尾（除非是根路径）
	if prefix != "" && prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix != "" && prefix != "/" && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	return &PrefixRouter{
		router: router,
		prefix: prefix,
	}
}

// joinPath 拼接前缀和路径
func (r *PrefixRouter) joinPath(path string) string {
	if r.prefix == "" || r.prefix == "/" {
		return path
	}
	if path == "/" {
		return r.prefix
	}
	return r.prefix + path
}

// GET 注册 GET 路由
func (r *PrefixRouter) GET(path string, handler http.HandlerFunc) {
	r.router.GET(r.joinPath(path), handler)
}

// POST 注册 POST 路由
func (r *PrefixRouter) POST(path string, handler http.HandlerFunc) {
	r.router.POST(r.joinPath(path), handler)
}

// Static 注册静态文件路由
func (r *PrefixRouter) Static(path, dir string) {
	r.router.Static(r.joinPath(path), dir)
}

// HandleFunc 注册任意 HTTP 方法的路由
func (r *PrefixRouter) HandleFunc(path string, handler http.HandlerFunc) {
	r.router.HandleFunc(r.joinPath(path), handler)
}

// SetPrefix 设置路由前缀
func (r *PrefixRouter) SetPrefix(prefix string) {
	if prefix != "" && prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix != "" && prefix != "/" && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	r.prefix = prefix
}

// GetPrefix 获取当前路由前缀
func (r *PrefixRouter) GetPrefix() string {
	return r.prefix
}

