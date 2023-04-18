package neo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// 原始的请求和响应对象
	Writer http.ResponseWriter
	Req    *http.Request

	// 当此请求方式
	Method string
	// 当此请求地址
	URL string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		URL:    r.URL.Path,
	}
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Status 设置响应状态码
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// HTML 返回HTML格式数据
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html)) // 不用处理
}

// JSON 返回JSON格式树
// JSON格式数据特殊点，需要给它先序列化
func (c *Context) JSON(code int, data interface{}) {
	c.SetHeader("Context-Type", "application/json")
	c.Status(code)
	// 序列化数据
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}

// String 返回纯文本格式数据
func (c *Context) String(code int, template string, value ...string) {
	c.SetHeader("Context-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(template, value)))
}

// Query 获取查询参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// TODO Param 获取请求参数，请求参数这里需要使用到动态路由再获取
//func (c *Context) Param(key string) string {}

// PostForm 获取请求体数据
// TODO 注意，根据用户传过来的数据格式的不同，获取数据的方式也是不同的。具体可以参考Gin
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}
