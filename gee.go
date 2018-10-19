package gee

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

// 控制器接口
// 方法、结构体都可以实现此方法
type GeeHandler interface {
	ServeHTTP(c *Context)
}

// 使func()实现GeeHandler
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
	Handler GeeHandler
	// 路由
	Route *Route
}

// 引擎
type Engine struct {
	Endpoints *[]Endpoint
	Router *httprouter.Router
	pool sync.Pool
}

// *** public methods ***

// 创建一个引擎
func CreateEngine() (Engine){
	return Engine{
		Endpoints: &([]Endpoint{}),
		Router: httprouter.New(),
		pool: sync.Pool{
			New: func () interface{} {
				return &Context{}
			},
		},
	}
}

// 注册路由
// 如果要修改结构体，这里e用地址传递
// FIXME err的传递
func (e *Engine) Register(method string, path string, handler GeeHandler) (err error) {
	route, err := createRoute(method, path)
	if err != nil {
		return err
	}
	handlers := append(*(e.Endpoints), Endpoint{
		Method: method,
		Path: path,
		Handler: handler,
		Route: route,
	})
	e.Endpoints = &handlers
	return

	// TODO 使用httprouter来匹配路由
	// e.Router.Handle(method, path, handler)
}

// 注册路由（方法）
func (e *Engine) RegisterFunc(method string, path string, fn func(*Context)) {
	e.Register(method, path, adaptFuncToGeeHandler(fn))
}

// 指定方法注册路由
// 调用的方法中修改了结构体，一样要用地址传递
func (e *Engine) GET(path string, handlerFunc func(*Context)) {
	e.RegisterFunc(http.MethodGet, path, handlerFunc)
}
func (e *Engine) POST(path string, handlerFunc func(*Context)) {
	e.RegisterFunc(http.MethodPost, path, handlerFunc)
}

// 中间件注册
func (e *Engine) USE(path string, handlerFunc func(*Context)) {
	e.RegisterFunc("ALL", path, handlerFunc)
}

// 处理请求
func (e Engine) HandleRequest(c *Context) {
	// 路径解析
	for _, ep := range *e.Endpoints {
		if ep.Route.Match(c) {
		// if c.Request.URL.Path == ep.Path && c.Request.Method == ep.Method {
			// 完全按照注册顺序排序
			c.HandlersChain = append(c.HandlersChain, ep.Handler)
			// ep.Handler.ServeHTTP(c)
		}
	}
	// 404处理
	if len(c.HandlersChain) == 0 {
		c.ResponseWriter.WriteHeader(404)
		c.ResponseWriter.Write([]byte("Not Found"))
		return
	}
	// 按注册顺序执行handlers
	handler := c.getNextHandler()
	handler.ServeHTTP(c)
}

// 实现http.Handler，将RequestWriter和Request转换为Context
func (e Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := e.pool.Get().(*Context)
	c.ResponseWriter = rw
	c.Request = req
	c.init()
	e.HandleRequest(c)
	e.pool.Put(c)
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