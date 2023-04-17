## HTTP 基础

### 标准库启动Web服务

在Go语言的`net/http`包启动一个Web服务式非常简单的，而且性能比别的语言性能要高很多，代码非常简单易懂

方式一：

```go
package main

import "net/http"

func main() {
	// 将路由和视图函数进行绑定
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("user home"))
	})

	// 启动服务器，第二个参数是nil表示用net/http内置的引擎处理器
	_ = http.ListenAndServe(":8080", nil)
}
```

上述的示例代码只是`net/http`启动Web服务的一种方式。

方式二：

```go
package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("user home"))
	})
	http.Handle("/user", mux)
	_ = http.ListenAndServe(":8080", mux)
}
```

对于这种启动方式，和方式一没有什么区别，只不过方式二使用的是自建的`ServeMux`实例对象，也就是之前说的IO多路复用器
细究一下，`ServeMux`结构体其实也是实现了`Handler`接口才能这样启动一个Web服务【再深入一下，就是该结构体实现了`ServeHTTP`
方法】。具体源码如下

```go
type Handler interface {
ServeHTTP(ResponseWriter, *Request)
}
================================ =
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
if r.RequestURI == "*" {
if r.ProtoAtLeast(1, 1) {
w.Header().Set("Connection", "close")
}
w.WriteHeader(StatusBadRequest)
return
}
h, _ := mux.Handler(r)
h.ServeHTTP(w, r)
}
```

方式三：
由第二个方式得出的灵感，我们是不是也可以自己实现`Handler`接口，然后用自定义的结构体充当IO多路复用器呢？

```go
package main

import (
	"fmt"
	"net/http"
)

type Engine struct{}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	switch url {
	case "/user":
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("hello %s", url)))
	case "/order":
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("hello %s", url)))
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("NOT FOUND"))
		return
	}
}

func NewEngine() *Engine {
	return &Engine{}
}

func main() {
	engine := NewEngine()
	_ = http.ListenAndServe(":8080", engine)
}

```

## neo框架雏形

在具体实现neo框架雏形的时候，我们需要先思考几个问题。可以对比Gin框架

1. 路由存储在什么地方
2. 路由怎么匹配
3. 路由怎么注册
4. Web服务如何启动
5. Web服务支持那些请求方式

```go
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

```