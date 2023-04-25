package main

import (
	"fmt"
	"github.com/borntodie-new/neo-web/neo/neo"
	"html/template"
	"net/http"
)

func withTemplateEngine(t neo.TemplateEngine) neo.TemplateOption {
	return func(engine *neo.Engine) {
		engine.T = t
	}
}

func main() {
	engine := neo.New()
	tpl, err := template.ParseGlob("./day06/template/*.gohtml")
	if err != nil {
		panic("Web 解析模板失败")
	}
	goTemplateEngine := neo.NewGoTemplateEngine(tpl)
	neo.WithTemplateOnEngine(engine, withTemplateEngine(goTemplateEngine))
	engine.GET("/login", func(ctx *neo.Context) {
		ctx.HTML(http.StatusOK, "login.gohtml", nil)
	})

	// 测试静态文件
	prefix := "file"
	s := neo.NewStaticFile("./day06/static", prefix)
	engine.GET(fmt.Sprintf("/assets/:%s", prefix), s.Handler())
	_ = engine.Run(":8080")
}
