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
	router       *router        // 路由树
	*RouterGroup                // 路由组
	groups       []*RouterGroup // 集中保存所有的路由组数据
}

// 对外对接用户，对内对接Web框架
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建Context上下文
	ctx := NewContext(w, r)
	// 方式一：获取所有的中间件方法，其实不可行，e.middlewares只能拿到Engine上的中间件
	// ctx.handlers = e.middlewares // 将该路由组的中间件函数传给上下文
	// 方式二：因为每个路由组身上都有可能有中间件，所以需要集中保存路由组
	// 但是不是所有的路由组的中间件列表都需要添加到当前上下文中
	// 只有那些符合的情况的才才需要：前缀一样的
	// 注意：这里只是添加中间件，命中的视图函数不是在这里添加
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			ctx.handlers = append(ctx.handlers, group.middlewares...)
		}
	}
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
	engine := &Engine{router: route, RouterGroup: group, groups: []*RouterGroup{}}
	group.engine = engine
	return engine
}

type RouterGroup struct {
	prefix      string        // 前缀
	parent      *RouterGroup  // 父路由组
	engine      *Engine       // 保存Engine对象完全是为了拿到Engine中的路由树信息，并对其注册和匹配
	middlewares []HandlerFunc // 存放中间件函数和视图函数。两者有优先级关系。中间件 > 视图函数
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix:      fmt.Sprintf("%s%s", group.prefix, prefix),
		parent:      group,
		engine:      group.engine,
		middlewares: group.middlewares,
	}
	// 每生成一个路由组，就需要添加到路由组列表中，方便集中管理
	group.engine.groups = append(group.engine.groups, newGroup)
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

// Use 中间件注册
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
