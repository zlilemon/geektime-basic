package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	exeName, _ := os.Executable()
	println("exeName:", exeName)
	server := gin.Default()
	server.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, go , hello")
		//c.JSON(200, gin.H{
		//	"message": "pong",
		//})
	})

	server.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "hello post")
	})

	server.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "这是参数路由"+name)

	})

	server.GET("/views/*.html", func(c *gin.Context) {
		name := c.Param(".html")
		c.String(http.StatusOK, "这是通配符路由"+name)

	})

	go func() {
		server1 := gin.Default()
		server1.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, "hello, go , hello, go , again")
			//c.JSON(200, gin.H{
			//	"message": "pong",
			//})
		})
		server1.Run(":8081")
	}()

	server.Run(":8080")
}
