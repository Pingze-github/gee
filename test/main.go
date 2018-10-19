package main

import (
	"fmt"
	"github.com/Pingze-github/gee"
	"net/http"
	"time"
)

type HomeController struct{}

func (h HomeController) ServeHTTP(c *gee.Context) {
	time.Sleep(time.Duration(1))
	c.End("Gee~ by Handler")
}

func foo(c *gee.Context) {
	time.Sleep(time.Duration(1))
	c.Write([]byte("Gee~ by HandlerFunc"))
	c.Next()
}

func timeoutHandler(c *gee.Context) {
	time.Sleep(time.Duration(3e9))
	c.Write([]byte("Done"))
}

func main() {
	homeController := HomeController{}

	//r := gin.Default()
	//r.GET("/", func(context *gin.Context) {
	//	time.Sleep(time.Duration(1))
	//	context.Done()
	//})
	//r.Run(":8080")

	engine := gee.CreateEngine()

	// 请求日志
	engine.USE("*", func(c *gee.Context) {
		timeStart := time.Now();
		defer func () {
			fmt.Println(fmt.Sprintf("[gee] %s %s %s", c.Method, c.Url, time.Since(timeStart)))
		}()
		fmt.Println(fmt.Sprintf("[gee] %s %s", c.Method, c.Url))
		c.Yield()
	})

	engine.Register(http.MethodGet, "/", homeController)


	// 双重定义
	engine.GET("/foo", foo)
	engine.GET("/foo", foo)


	engine.GET("/timeout", timeoutHandler)
	// engine.Get("/hello", homeController)
	err := engine.Start("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
}