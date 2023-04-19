package neo

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // 路由树，其实应该是路由森林，每一个HTTP method都有一颗树【key是method，value是节点】
	handlers map[string]HandlerFunc // 路由和视图作绑定【key是路由，value是视图函数】
}

// 内部核心API，仅共内部使用，用于注册路由
// 注册路由
// method = GET | pattern = "/user/home" => parts = [user, home]
// method = GET | pattern = "/user/order" => parts = [user, order]
func (r *router) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	root, ok := r.roots[method]
	if !ok { // 根路由树不存在
		root = &node{
			part: "/",
		}
		r.roots[method] = root
	}
	// pattern 必须以 / 开头
	if !strings.HasPrefix(pattern, "/") {
		panic("web: 路由必须以 / 开头")
	}
	if strings.HasSuffix(pattern, "/") {
		panic("web: 路由不能以 / 结尾")
	}
	// 切割pattern
	parts := strings.Split(pattern[1:], "/")
	for _, part := range parts {
		if part == "" {
			panic("web: 路由不能连续出现 / ")
		}
		child := root.search(part)
		if child == nil { // 返回值是nil，说明当前节点下没有part的节点，需要新增
			child = &node{ // 这里不需要添加pattern字段，只有最底层的叶子节点才需要添加。我们统一在循环之后再添加
				part: part,
			}
			// 如果part中带有 : 表示是模糊匹配，需要给节点的isWild设为true
			if strings.HasPrefix(part, ":") {
				child.isWild = true
			}
			root.children = append(root.children, child)
		}
		root = child // 查找到了child，沿着child继续查找
	}
	root.pattern = pattern
	key := fmt.Sprintf("%s-%s", method, pattern)
	r.handlers[key] = handlerFunc
	log.Printf("Add Router %4s - %s", method, pattern)
}

// 匹配路由
func (r *router) getRouter(method string, pattern string) (*node, map[string]string) {
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		// 路由树都不存在，直接返回nil
		return nil, nil
	}
	// 切割pattern
	pattern = strings.Trim(pattern, "/")
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if part == "" {
			return nil, nil
		}
		child := root.search(part)
		if child == nil {
			return nil, nil
		}
		// 将匹配到的带有:的路由添加到params中
		if child.isWild {
			params[part[1:]] = part
		}
		// 1. child.pattern != ""这是精确匹配，但是还不够，有点缺陷
		// 2. child.pattern == fmt.Sprintf("/%s", pattern) 这是精确匹配，搭配条件1才是完美
		// 上述两个条件只有配合使用才能精确匹配路由。/login 和 /login/123 两种路由都能匹配到
		// 3. child.isWild 这是模糊匹配，优先级低于精确匹配
		if (child.pattern != "" && child.pattern == fmt.Sprintf("/%s", pattern)) || child.isWild {
			return child, params
		}
		root = child
	}
	return nil, nil
}

func (r *router) handle(ctx *Context) {
	// 请求来了，需要匹配路由
	log.Printf("Request %4s - %s", ctx.Method, ctx.URL)
	n := r.getRouter(ctx.Method, ctx.URL)
	if n == nil {
		// 没有匹配到
		ctx.String(http.StatusInternalServerError, "NOT FOUND")
		return
	}
	key := fmt.Sprintf("%s-%s", ctx.Method, n.pattern)
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
	return &router{roots: map[string]*node{}, handlers: map[string]HandlerFunc{}}
}
