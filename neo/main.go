package main

import (
	"fmt"
	"github.com/borntodie-new/neo-web/neo/neo"
	"net/http"
)

func main() {
	engine := neo.New()
	engine.GET("/user", func(ctx *neo.Context) {
		username := ctx.Query("username")
		password := ctx.Query("password")
		ctx.HTML(http.StatusOK, fmt.Sprintf("<h1>欢迎：%s</h1><br/><h1>你的密码是：%s<h1>", username, password))
	})
	engine.POST("/login", func(ctx *neo.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		ctx.JSON(http.StatusOK, struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: username,
			Password: password,
		})
	})
	_ = engine.Run(":8080")
}
