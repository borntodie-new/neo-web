package main

import (
	"neo-web/neo/neo"
	"net/http"
)

func main() {
	engine := neo.New()
	engine.GET("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("user"))
	})
	engine.POST("/order", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("order"))
	})
	_ = engine.Run(":8080")
}
