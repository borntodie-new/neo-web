package main

import (
	"fmt"
	"github.com/borntodie-new/neo-web/day06/neo"
	"html/template"
	"net/http"
)

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
	// 测试模板渲染
	engine.GET("/login", func(ctx *neo.Context) {
		ctx.HTML(http.StatusOK, "login.gohtml", nil)
	})

	// 测试静态文件
	prefix := "file"
	s := neo.NewStaticFile("./day06/static", prefix)
	engine.GET(fmt.Sprintf("/assets/:%s", prefix), s.Handler())
	_ = engine.Run(":8080")
}
