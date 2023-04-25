package main

import (
	"fmt"
	"github.com/borntodie-new/neo-web/day07/neo"
	"net/http"
)

func main() {
	engine := neo.Default()
	engine.GET("/user", func(ctx *neo.Context) {
		data := []string{"A", "B", "C"}
		fmt.Println(data[1000]) // 肯定报错，下标越界
	})
	engine.GET("/order", func(ctx *neo.Context) {
		fmt.Println("order成功")
		ctx.String(http.StatusOK, "order成功")
	})
	_ = engine.Run(":8080")
}
