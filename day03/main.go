package main

import (
	"github.com/borntodie-new/neo-web/day03/neo"
	"net/http"
)

func main() {
	engine := neo.New()
	engine.GET("/", func(ctx *neo.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Neo</h1>")
	})
	engine.GET("/hello", func(ctx *neo.Context) {
		// expect /hello?name=neo
		ctx.JSON(http.StatusOK, neo.H{
			"code": "200",
			"name": ctx.Query("name"),
		})
	})
	engine.GET("/hello/:name", func(ctx *neo.Context) {
		// expect /hello/neo
		ctx.JSON(http.StatusOK, neo.H{
			"code": "200",
			"name": ctx.Params("name"),
		})
	})
	engine.GET("/assets/*filepath", func(ctx *neo.Context) {
		// expect /hello/neo
		ctx.JSON(http.StatusOK, neo.H{
			"code": "200",
			"name": ctx.Params("filepath"),
		})
	})
	_ = engine.Run(":8080")
}
