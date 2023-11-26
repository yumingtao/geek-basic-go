package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 一个web服务器被抽象为engine
	// engine承担了路由注册，接入middleware的核心职责
	server := gin.Default()
	server.Use(func(context *gin.Context) {
		println("第一个middleware")
	}, func(context *gin.Context) {
		println("第二个middleware")
	})
	// 路由注册，注册了一个get方法
	// gin支持的路由类型：静态路由，完全匹配的路由；参数路由：在路径中带上参数的路由；通配符路由：任意匹配的路由
	server.GET("/hello", func(context *gin.Context) {
		// context核心职责：处理请求，返回响应
		context.String(http.StatusOK, "Hello, World!")
	})

	// 参数路由, 路径参数
	server.GET("/users/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "Hello,"+name)
	})

	// 参数路由, 查询参数 /order?id
	server.GET("/orders", func(context *gin.Context) {
		id := context.Query("id")
		context.String(http.StatusOK, "Your order is "+id)
	})

	// 通配符路由, 注意在gin中*不能单独出现
	// /views/* 和 /views/*/id 都不行
	server.GET("/views/*.html", func(context *gin.Context) {
		view := context.Param(".html")
		context.String(http.StatusOK, "view is "+view)
	})

	// 可以创建并启动多个engine
	/*go func() {
		server1 := gin.Default()
		server1.Run(":8081")
	}()*/
	// 如果不传参数，默认监听的是8080端口
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
