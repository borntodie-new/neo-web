package neo

import (
	"fmt"
	"net/http"
)

// HandlerFunc 视图函数签名
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc // 路由和视图作绑定【key是路由，value是视图函数】
}

// 内部核心API，仅共内部使用，用于注册路由
func (e *Engine) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	key := fmt.Sprintf("%s-%s", method, pattern)
	e.router[key] = handlerFunc
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
	// 请求来了，需要匹配路由
	url := r.URL.Path
	method := r.Method
	key := fmt.Sprintf("%s-%s", method, url)
	handlerFunc, ok := e.router[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("NOT FOUND"))
		return
	}
	// 执行命中的视图函数
	handlerFunc(w, r)
}

// Run 手动启动服务，控制力强
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}
