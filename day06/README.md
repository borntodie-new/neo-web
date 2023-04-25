# 模板渲染——template
目前大部分系统、网站都是采用前后端分离的架构了。就是前端只负责画页面和请求数据，后端只负责响应请求。
一般后端响应的数据都是`JSON`和`ProtocolBuffer`格式

而对于模板渲染，这是属于前后端不分离的架构使用的。意思就是前端代码和后端代码都是由后端程序员完成，整个前端页面都是在服务端渲染好了，然后
直接将这些模板返回给前端。


在模板渲染这一块有两个内容
### 静态文件的访问
关于静态文件，我们这里只是实现一个非常非常简单的版本，只是为了完善咱们框架的功能而已，不能用于生产环境中，因为性能太差，如果有条件，直接使用OSS和CDN服务就好。

我这里还想说一下，静态文件有些需要注意事项
1. 静态文件是否需要缓存
2. 如何缓存，或者说用那种算法缓存
3. 缓存文件的大小是否有限制
4. 缓存文件的数量是否有限制
5. ....

具体的实现有我认为有两种方式
1. 使用`net/http`包中提供了一个文件处理的结构体`http.FileServer`，具体实现是`serveFile`方法，它是一个内置的方法。具体原理大家可以看代码实现
2. 自己实现一个文件读取返回给前端

其实兔兔就是使用的是第一种方式，对于第一种方式，性能不是很好，当然我们的方式性能也没好到哪里去，只不过为了给大家讲一下静态文件是怎么读取的并且返回给前端

我们是自己实现一个文件读取返回给前端

具体实现
```go
type StaticFile struct {
	Dir  string // 需要开放的文件路径
	Path string // 参数地址
}

func NewStaticFile(dir string, path string) *StaticFile {
	return &StaticFile{Dir: dir, Path: path}
}

func (s *StaticFile) Handler() HandlerFunc {
	return func(ctx *Context) {
		// 拿到URL中的文件名
		fileName := ctx.Params(s.Path)
		// 拼接文件地址并打开
		filePath := filepath.Join(s.Dir, fileName)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "服务器异常")
			return
		}
		ctx.String(http.StatusOK, string(data))
	}
}

```

具体使用
```go
func main() {
	engine := neo.New()
	// 测试静态文件
	prefix := "file"
	s := neo.NewStaticFile("./day06/static", prefix)
	engine.GET(fmt.Sprintf("/assets/:%s", prefix), s.Handler())
	_ = engine.Run(":8080")
}
```

### HTML模板的渲染
对于HTML模板，我们这是直接抽象出一个接口，用于对接外部的模板语言，而Web框架只是对接我们的模板接口

而且对于这个功能，应该不算是框架的核心的功能。所以提供一个可选的方法给用户使用，**我们必须记住，不能让多数人给少数人买单**。
我们做开源框架的时候，一定要牢记这个：只提供一些方法，并且这个方法还是比较通用的，而对于一个用户特殊的需求是需要用户自己根据框架提供的方法组装实现自己的需求。

模板引擎抽象
```go
// TemplateEngine 模板引擎抽象
type TemplateEngine interface {
	// Render 渲染页面方法
	// ctx 上下文，可能需要从中拿取相应树
	// tplName 模板名字
	// data 需要填充到模板中的数据
	// 返回值 []byte渲染后的模板数据
	// 返回值 error错误信息
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}
```
具体实现——内置的模板
```go
// 使用Golang内置的模板来充当模板引擎
type GoTemplateEngine struct {
	T *template.Template
}

func NewGoTemplateEngine(t *template.Template) TemplateEngine {
	return &GoTemplateEngine{T: t}
}

// Render 渲染数据
func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	// 这里的任务就是将 data 渲染到 模板名是 tplName 的模板中。
	// 那就是用 html/template 包
	// 从哪来？或者说，怎么渲染
	buf := &bytes.Buffer{}
	// ExecuteTemplate：将data渲染到tplName中，并将最后出来的结果放在buf中
	err := g.T.ExecuteTemplate(buf, tplName, data)
	return buf.Bytes(), err
}
```
内嵌模板引擎到Server和Context上下文中
```go
type Engine struct {
	router       *router        // 路由树
	*RouterGroup                // 路由组
	groups       []*RouterGroup // 集中保存所有的路由组数据

	T TemplateEngine // 模板引擎
}

type Context struct {
    // 原始的请求和响应对象
    Writer http.ResponseWriter
    Req    *http.Request
    
    // 当此请求方式
    Method string
    // 当此请求地址
    URL string
    // 请求参数 不需要暴露出去
    params map[string]string
    
    // 中间件
    handlers []HandlerFunc
    index    int
    
    // 模板引擎
    t TemplateEngine
}
```
我们之前说了，模板语言这个功能只要部分用户需要，所以我们这里需要构建一个可配置的选项给到用户使用模板引擎。具体实现是用一个Option设计模式
```go
// Option模式设计
type EngineOption func(engine *Engine)

func WithEngineOptions(engine *Engine, opts ...EngineOption) {
	for _, opt := range opts {
		opt(engine)
	}
}

// Option模式使用
func withTemplate(t neo.TemplateEngine) neo.EngineOption {
    return func(engine *neo.Engine) {
        engine.T = t
    }
}

func main() {
    engine := neo.New()
    tpl, err := template.ParseGlob("day06/template/*.gohtml")
    if err != nil {
        panic("模板解析错误" + err.Error())
    }
    // 这里需要将模板引擎注册到
    goTemplateEngine := neo.NewGoTemplateEngine(tpl)
    neo.WithEngineOptions(engine, withTemplate(goTemplateEngine))
    engine.GET("/login", func(ctx *neo.Context) {
        ctx.HTML(http.StatusOK, "login.gohtml", nil)
    })
    _ = engine.Run(":8080")
}



```
关于Option设计模式，大家自行上网百度即可。

### 总结
我们要清楚，上述说的是两个功能
- 模板渲染
- 静态文件的开放

我们是将这两个功能分开实现的。兔兔文章里好像没有直接说明这点，所以看起来其实是有点懵的