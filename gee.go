package gee

import (
	"fmt"
	"net/http"
)

type Endpoint struct {
	// http方法
	Verb string
	// 路径
	Path string
	// 控制器
	Handler GeeHandler
	HandlerFunc func(*Context)
	isFunc bool
}

// 引擎，内部包含很多Endpoint
type Engine struct {
	Endpoints *[]Endpoint
}

// 创建一个引擎
func CreateEngine() (Engine){
	return Engine{Endpoints: &( []Endpoint{} )}
}

// 注册路由
// 如果要修改结构体，这里e用地址传递
func (e *Engine) Register(verb string, path string, handler GeeHandler) {
	handlers := append(*(e.Endpoints), Endpoint{
		Verb: verb,
		Path: path,
		Handler: handler,
	})
	e.Endpoints = &handlers
}

func (e *Engine) RegisterFunc(verb string, path string, fn func(*Context)) {
	handlers := append(*(e.Endpoints), Endpoint{
		Verb: verb,
		Path: path,
		HandlerFunc: fn,
		isFunc: true,
	})
	e.Endpoints = &handlers
}

// 指定方法注册路由
// 调用的方法中修改了结构体，一样要用地址传递
func (e *Engine) Get(path string, handlerFunc func(*Context)) {
	e.RegisterFunc(http.MethodGet, path, handlerFunc)
}
func (e *Engine) Post(path string, handlerFunc func(*Context)) {
	e.RegisterFunc(http.MethodPost, path, handlerFunc)
}


// 上下文
type Context struct {
	ResponseWriter http.ResponseWriter
	Request *http.Request
}

// 控制器接口
type GeeHandler interface {
	ServeHTTP(c *Context)
}

// 处理请求
func (c *Context) Handle(e Engine) {
	// 路径解析
	for _, ep := range *e.Endpoints {
		if c.Request.URL.Path == ep.Path && c.Request.Method == ep.Verb {
			if ep.isFunc == true {
				ep.HandlerFunc(c)
			} else {
				ep.Handler.ServeHTTP(c)
			}
			return
		}
	}
	c.ResponseWriter.WriteHeader(404)
	c.ResponseWriter.Write([]byte("Not Found"))
	return
}

// 实现http.Handler，将RequestWriter和Request转换为Context
func (e Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := Context{rw, req}
	c.Handle(e)
}

// 启动服务
func (e Engine) Start(addr string) (error) {
	// r := mux.NewRouter()
	for _, v := range *e.Endpoints {
		fmt.Println(fmt.Sprintf("Server Register route: %s %s", v.Verb, v.Path))
		// r.Handle(v.Path, v.Handler).Methods(v.Verb)
	}
	// http.Handle("/", r)
	http.Handle("/", e)
	fmt.Println(fmt.Sprintf("Server running @ http://%s", addr))
	err := http.ListenAndServe(addr, nil)
	return err
}