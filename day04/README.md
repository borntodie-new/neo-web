# 路由分组
关于路由分组的好处，我这就不细说了，对于这个功能，其实可以没有的，只是为了方便用户使用，我们额外加上这个功能

具体看下Gin框架的路由风分组如何使用的
```go
package main

import "github.com/gin-gonic/gin"

func main() {
   r := gin.Default()
   r.Group("v1")
   {
	   r.GET("/user", func(ctx *gin.Context) {})
	   r.GET("/order", func(ctx *gin.Context) {})
   }
   // http://localhost:8080/v1/user
   // http://localhost:8080/v1/order
   r.Run(":8080")
}
```

## 我们想实现的效果
1. engine能够能够调用`Group`方法实现分组
2. `GET`、`POST`、`PUT`、`DELETE`等衍生API转移到`RouterGroup`中实现
3. `addRouter`方法转移到`RouterGroup`中实现，统一和路由树沟通

通过上述的分析，得出结论，`RouterGroup`结构体字段如下
```go
type RouterGroup struct {
	prefix string       // 前缀
	parent *RouterGroup // 父路由组
	engine *Engine
}
```
通过上述的分析，得出结论，`Engine`结构更新如下
```go
type Engine struct {
	router       *router // 路由树
	*RouterGroup         // 路由组
}
```
完成效果1
```go
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: fmt.Sprintf("%s%s", group.prefix, prefix),
		parent: group,
		engine: group.engine,
	}
	return newGroup
}
```
完成效果2
```go
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
```
完成效果3
```go
// 在路由组上定义一个添加路由的方法，这个作为唯一和路由树交互的入口
func (group *RouterGroup) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	pattern = fmt.Sprintf("%s%s", group.prefix, pattern)
	log.Printf("Add Router %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handlerFunc)
}
```
上述完成测试
```go
func main() {
    engine := neo.New()
    v1 := engine.Group("/v1")
    {
        v1.GET("/user", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/user")
        })
        v1.POST("/order", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/order")
        })
        v1.DELETE("/cart", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/cart")
        })
        v1.PUT("/admin", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/admin")
        })
    }
    v2 := engine.Group("/v2")
    {
        v2.GET("/user", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/user")
        })
        v2.POST("/order", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/order")
        })
        v2.DELETE("/cart", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/cart")
        })
        v2.PUT("/admin", func(ctx *neo.Context) {
            ctx.String(http.StatusOK, "v1/admin")
        })
    }
    _ = engine.Run(":8080")
}
```