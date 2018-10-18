package main

import (
	"github.com/Pingze-github/gee"
	"net/http"
	"time"
)

type HomeController struct{}

func (h HomeController) ServeHTTP(c *gee.Context) {
	c.ResponseWriter.Write([]byte("Gee~ by Handler"))
}

func foo(c *gee.Context) {
	c.ResponseWriter.Write([]byte("Gee~ by HandlerFunc"))
}

func timeoutHandler(c *gee.Context) {
	time.Sleep(time.Duration(3e9))
	c.ResponseWriter.Write([]byte("Done"))
}

func main() {
	homeController := HomeController{}

	engine := gee.CreateEngine()
	engine.Register(http.MethodGet, "/", homeController)
	engine.Get("/func", foo)
	engine.Get("/timeout", timeoutHandler)
	// engine.Get("/hello", homeController)
	engine.Start("127.0.0.1:8080")
}