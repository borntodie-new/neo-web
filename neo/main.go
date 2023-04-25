package main

import (
	"fmt"
	"github.com/borntodie-new/neo-web/neo/neo"
	"net/http"
)

func withTemplateEngine(t neo.TemplateEngine) neo.TemplateOption {
	return func(engine *neo.Engine) {
		engine.T = t
	}
}

func main() {
	engine := neo.Default()
	engine.GET("/user", func(ctx *neo.Context) {
		data := []string{"A", "B", "C"}
		fmt.Println(data[1000]) // 绝对报错
	})
	engine.GET("/order", func(ctx *neo.Context) {
		ctx.String(http.StatusOK, "order请求成功")
	})
	_ = engine.Run(":8080")
}
