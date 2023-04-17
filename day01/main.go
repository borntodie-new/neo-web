package main

import (
	"fmt"
	"net/http"
)

// 实现方式一：使用net/http默认的IO多路复用器
//func main() {
//	// 将路由和视图函数进行绑定
//	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		_, _ = w.Write([]byte("user home"))
//	})
//
//	// 启动服务器，第二个参数是nil表示用net/http内置的IO多路复用器
//	_ = http.ListenAndServe(":8080", nil)
//}

// 实现方式二：使用net/http实现好的ServeMux结构体充当IO多路复用器
//func main() {
//	mux := http.NewServeMux()
//	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		_, _ = w.Write([]byte("user home"))
//	})
//	http.Handle("/user", mux)
//	_ = http.ListenAndServe(":8080", mux)
//}

// 方式三：自定义结构体实现Handler接口，充当IO多路复用器
func main() {
	engine := NewEngine()
	_ = http.ListenAndServe(":8080", engine)
}

type Engine struct{}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	switch url {
	case "/user":
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("hello %s", url)))
	case "/order":
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("hello %s", url)))
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("NOT FOUND"))
		return
	}
}

func NewEngine() *Engine {
	return &Engine{}
}
