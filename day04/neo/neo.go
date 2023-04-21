package neo

import (
	"fmt"
	"log"
	"net/http"
)

type H map[string]string

// HandlerFunc 视图函数签名
type HandlerFunc func(ctx *Context)

type Engine struct {
	router       *router // 路由树
	*RouterGroup         // 路由组
}

// 对外对接用户，对内对接Web框架
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建Context上下文
	ctx := NewContext(w, r)
	// 转发请求到框架
	e.router.handle(ctx)
}

// Run 手动启动服务，控制力强
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	route := newRouter()
	group := &RouterGroup{}
	engine := &Engine{router: route, RouterGroup: group}
	group.engine = engine
	return engine
}

type RouterGroup struct {
	prefix string       // 前缀
	parent *RouterGroup // 父路由组
	engine *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: fmt.Sprintf("%s%s", group.prefix, prefix),
		parent: group,
		engine: group.engine,
	}
	return newGroup
}

// 在路由组上定义一个添加路由的方法，这个作为唯一和路由树交互的入口
func (group *RouterGroup) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = fmt.Sprintf("%s%s", group.prefix, pattern)
	log.Printf("Add Router %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handlerFunc)
}

// GET 外部衍生API，提供给用户使用，现在嫁接到RouterGroup上
func (group *RouterGroup) GET(pattern string, handlerFunc HandlerFunc) {
	group.addRouter(http.MethodGet, pattern, handlerFunc)
}
func (group *RouterGroup) POST(pattern string, handlerFunc HandlerFunc) {
	group.addRouter(http.MethodPost, pattern, handlerFunc)
}
func (group *RouterGroup) DELETE(pattern string, handlerFunc HandlerFunc) {
	group.addRouter(http.MethodDelete, pattern, handlerFunc)
}
func (group *RouterGroup) PUT(pattern string, handlerFunc HandlerFunc) {
	group.addRouter(http.MethodPut, pattern, handlerFunc)
}
