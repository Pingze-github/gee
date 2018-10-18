package gee

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// 控制器接口
// 方法、结构体都可以实现此方法
type GeeHandlerInterface interface {
	ServeHTTP(c *Context)
}

// 使func()实现GeeHandlerInterface接口
type adaptFuncToGeeHandler func(*Context)

func (f adaptFuncToGeeHandler) ServeHTTP(c *Context) {
	f(c)
}

// 路由节点
type Endpoint struct {
	// http方法
	Method string
	// 路径
	Path string
	// 控制器
	Handler GeeHandlerInterface
}

// 引擎
type Engine struct {
	Endpoints *[]Endpoint
	Router *httprouter.Router
}

// 上下文
type Context struct {
	ResponseWriter http.ResponseWriter
	Request *http.Request
}

// *** public methods ***

// 创建一个引擎
func CreateEngine() (Engine){
	return Engine{
		Endpoints: &([]Endpoint{}),
		Router: httprouter.New(),
	}
}

// 注册路由
// 如果要修改结构体，这里e用地址传递
func (e *Engine) Register(method string, path string, handler GeeHandlerInterface) {
	handlers := append(*(e.Endpoints), Endpoint{
		Method: method,
		Path: path,
		Handler: handler,
	})
	e.Endpoints = &handlers

	// TODO 使用httprouter来匹配路由
	// e.Router.Handle(method, path, handler)
}

// 注册路由（方法）
func (e *Engine) RegisterFunc(method string, path string, fn func(*Context)) {
	handlers := append(*(e.Endpoints), Endpoint{
		Method: method,
		Path: path,
		Handler: adaptFuncToGeeHandler(fn),
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

// 处理请求
func (c *Context) Handle(e *Engine) {
	// 路径解析
	for _, ep := range *e.Endpoints {
		if c.Request.URL.Path == ep.Path && c.Request.Method == ep.Method {
			ep.Handler.ServeHTTP(c)
			return
		}
	}
	c.ResponseWriter.WriteHeader(404)
	c.ResponseWriter.Write([]byte("Not Found"))
	return
}

// 实现http.Handler，将RequestWriter和Request转换为Context
func (e Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	go func() {
		c := Context{rw, req}
		c.Handle(&e)
	}()
}

// 启动服务
func (e Engine) Start(addr string) (error) {
	for _, v := range *e.Endpoints {
		fmt.Println(fmt.Sprintf("Server Register route: %s %s", v.Method, v.Path))
		// r.Handle(v.Path, v.Handler).Methods(v.Method)
	}
	// http.Handle("/", r)
	http.Handle("/", e)
	fmt.Println(fmt.Sprintf("Server running @ http://%s", addr))
	err := http.ListenAndServe(addr, nil)
	return err
}