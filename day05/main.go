package main

import (
	"fmt"
	"github.com/borntodie-new/neo-web/day05/neo"
	"net/http"
)

func Logger1() neo.HandlerFunc {
	return func(ctx *neo.Context) {
		fmt.Println("请求来了Logger1")
		ctx.Next()
		fmt.Println("请求走了Logger1")
	}
}
func Logger2() neo.HandlerFunc {
	return func(ctx *neo.Context) {
		fmt.Println("请求来了Logger2")
		ctx.JSON(200, neo.H{
			"code": "200",
			"msg":  "Logger2退出",
		})
		ctx.Abort()
		//ctx.Next()
		fmt.Println("请求走了Logger2")
	}
}
func Logger3() neo.HandlerFunc {
	return func(ctx *neo.Context) {
		fmt.Println("请求来了Logger3")
		ctx.Next()
		fmt.Println("请求走了Logger3")
	}
}
func Logger4() neo.HandlerFunc {
	return func(ctx *neo.Context) {
		fmt.Println("请求来了Logger4")
		ctx.Next()
		fmt.Println("请求走了Logger4")
	}
}
func Logger5() neo.HandlerFunc {
	return func(ctx *neo.Context) {
		fmt.Println("请求来了Logger5")
		ctx.Next()
		fmt.Println("请求走了Logger5")
	}
}
func main() {
	engine := neo.New()
	engine.Use(Logger1())
	v1 := engine.Group("/v1")
	v1.Use(Logger2())
	{
		v1.GET("/user", func(ctx *neo.Context) {
			fmt.Println("v1/user")
			ctx.String(http.StatusOK, "v1/user")
		})
		v1.POST("/order", func(ctx *neo.Context) {
			fmt.Println("v1/order")
			ctx.String(http.StatusOK, "v1/order")
		})
		v1.DELETE("/cart", func(ctx *neo.Context) {
			fmt.Println("v1/cart")
			ctx.String(http.StatusOK, "v1/cart")
		})
		v1.PUT("/admin", func(ctx *neo.Context) {
			fmt.Println("v1/admin")
			ctx.String(http.StatusOK, "v1/admin")
		})
	}
	v2 := engine.Group("/v2")
	v2.Use(Logger3(), Logger4(), Logger5())
	{
		v2.GET("/user", func(ctx *neo.Context) {
			fmt.Println("/v2/user")
			ctx.String(http.StatusOK, "v2/user")
		})
		v2.POST("/order", func(ctx *neo.Context) {
			fmt.Println("/v2/order")
			ctx.String(http.StatusOK, "v2/order")
		})
		v2.DELETE("/cart", func(ctx *neo.Context) {
			fmt.Println("/v2/cart")
			ctx.String(http.StatusOK, "v2/cart")
		})
		v2.PUT("/admin", func(ctx *neo.Context) {
			fmt.Println("/v2/admin")
			ctx.String(http.StatusOK, "v2/admin")
		})
	}
	_ = engine.Run(":8080")
}
