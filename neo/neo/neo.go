package neo

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type H map[string]string

// HandlerFunc 视图函数签名
type HandlerFunc func(ctx *Context)

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup // 保存所有的路由组信息，方便后期匹配路由组
}

// 对外对接用户，对内对接Web框架
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建Context上下文
	ctx := NewContext(w, r)
	// 请求来的时候，收集匹配当前URL地址的中间件函数
	// 注意，这里只匹配中间件
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			ctx.handlers = append(ctx.handlers, group.middlewares...)
		}
	}
	// 转发请求到框架
	// 里面匹配命中的视图函数
	e.router.handle(ctx)
}

// Run 手动启动服务，控制力强
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	r := newRouter()
	routerGroup := &RouterGroup{}
	engine := &Engine{
		router:      r,
		RouterGroup: routerGroup,
		groups:      []*RouterGroup{},
	}
	routerGroup.engine = engine
	return engine
}

type RouterGroup struct {
	prefix      string        // 路由组前缀
	parent      *RouterGroup  // 父级路由组
	engine      *Engine       // 完全是为了路由组能够拿到路由树，而路由树又在Engine中
	middlewares []HandlerFunc // 当前路由组中注册的所有中间件函数
}

func (group *RouterGroup) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = fmt.Sprintf("%s%s", group.prefix, pattern)
	group.engine.router.addRouter(method, pattern, handlerFunc)
	log.Printf("Add Router %4s - %s", method, pattern)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// 处理用户没有以 / 开头
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}
	newGroup := &RouterGroup{
		prefix:      fmt.Sprintf("%s%s", group.prefix, prefix),
		parent:      group,
		engine:      group.engine,
		middlewares: group.middlewares, // 必须填充父级的中间件方法列表
	}
	// 向Engine的路由组列表字段添加新建的路由组
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// GET 外部衍生API，提供给用户使用
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

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
