package neo

import (
	"fmt"
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc // 路由和视图作绑定【key是路由，value是视图函数】
}

// 内部核心API，仅共内部使用，用于注册路由
func (r *router) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	key := fmt.Sprintf("%s-%s", method, pattern)
	log.Printf("Add Router %4s - %s", method, pattern)
	r.handlers[key] = handlerFunc
}

func (r *router) handle(ctx *Context) {
	// 请求来了，需要匹配路由
	key := fmt.Sprintf("%s-%s", ctx.Method, ctx.URL)
	log.Printf("Request %4s - %s", ctx.Method, ctx.URL)
	handlerFunc, ok := r.handlers[key]
	if !ok {
		ctx.String(http.StatusInternalServerError, "NOT FOUND")
		return
	}
	// 构建Context请求上下文
	// 执行命中的视图函数
	handlerFunc(ctx)
}

func newRouter() *router {
	return &router{handlers: map[string]HandlerFunc{}}
}
