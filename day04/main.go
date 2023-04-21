package main

import (
	"github.com/borntodie-new/neo-web/day04/neo"
	"net/http"
)

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
