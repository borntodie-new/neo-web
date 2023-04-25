package neo

import (
	"bytes"
	"context"
	"html/template"
)

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

// ParseGlob 解析模板
//func (g *GoTemplateEngine) ParseGlob(tplName string) error {
//	var err error
//	g.T, err = g.T.ParseGlob(tplName)
//	return err
//}
