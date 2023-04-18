package neo

import (
	"net/http"
)

// HandlerFunc 视图函数签名
type HandlerFunc func(ctx *Context)

type Engine struct {
	*router
}

// GET 外部衍生API，提供给用户使用
func (e *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	e.addRouter(http.MethodGet, pattern, handlerFunc)
}
func (e *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	e.addRouter(http.MethodPost, pattern, handlerFunc)
}
func (e *Engine) DELETE(pattern string, handlerFunc HandlerFunc) {
	e.addRouter(http.MethodDelete, pattern, handlerFunc)
}
func (e *Engine) PUT(pattern string, handlerFunc HandlerFunc) {
	e.addRouter(http.MethodPut, pattern, handlerFunc)
}

// 对外对接用户，对内对接Web框架
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建Context上下文
	ctx := NewContext(w, r)
	// 转发请求到框架
	e.handle(ctx)
}

// Run 手动启动服务，控制力强
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	return &Engine{newRouter()}
}
