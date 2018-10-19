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
	Serve(c *Context)
}

// 使func()实现GeeHandler
type adaptFuncToGeeHandler func(*Context)

func (f adaptFuncToGeeHandler) Serve(c *Context) {
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

// 默认最终处理中间件
func finalHandler(c *Context, data interface{}) {}

// 引擎
type Engine struct {
	Endpoints *[]Endpoint
	Router *httprouter.Router
	pool sync.Pool
	FinalHandler func(*Context, interface{})
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
		FinalHandler: finalHandler,
	}
}

// 注册路由
// 如果要修改结构体，这里e用地址传递
// FIXME err的传递
func (e *Engine) Register(method string, path string, handlers ...GeeHandler) (err error) {
	route, err := createRoute(method, path)
	if err != nil {
		return err
	}
	var endpoints []Endpoint
	for _, handler := range(handlers) {
		endpoints = append(*e.Endpoints, Endpoint{
			Method: method,
			Path: path,
			Handler: handler,
			Route: route,
		})
		e.Endpoints = &endpoints
	}
	return

	// TODO 使用httprouter来匹配路由
	// e.Router.Handle(method, path, handler)
}

// 注册路由（方法）
func (e *Engine) RegisterFunc(method string, path string, fns ...func(*Context)) {
	var handlers []GeeHandler
	for _, fn := range(fns) {
		handlers = append(handlers, adaptFuncToGeeHandler(fn))
	}
	e.Register(method, path, handlers...)
}

// 指定方法注册路由
// 调用的方法中修改了结构体，一样要用地址传递
func (e *Engine) GET(path string, handlerFuncs ...func(*Context)) {
	e.RegisterFunc(http.MethodGet, path, handlerFuncs...)
}
func (e *Engine) POST(path string, handlerFuncs ...func(*Context)) {
	e.RegisterFunc(http.MethodPost, path, handlerFuncs...)
}

// 中间件注册
func (e *Engine) USE(path string, handlerFuncs ...func(*Context)) {
	e.RegisterFunc("ALL", path, handlerFuncs...)
}

// Final处理中间件注册
func (e *Engine) Final(finalFunc func(*Context, interface{})) {
	e.FinalHandler = finalFunc
}

// 处理请求
func (e *Engine) HandleRequest(c *Context) {
	// 路径解析
	for _, ep := range *e.Endpoints {
		if ep.Route.Match(c) {
		// if c.Request.URL.Path == ep.Path && c.Request.Method == ep.Method {
			// 完全按照注册顺序排序
			c.HandlersChain = append(c.HandlersChain, ep.Handler)
			// ep.Handler.ServeHTTP(c)
		}
	}
	// 按注册顺序执行handlers
	handler := c.getNextHandler()
	handler.Serve(c)
}

// 实现http.Handler接口
// 将RequestWriter和Request转换为Context
func (e *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := e.pool.Get().(*Context)
	c.ResponseWriter = rw
	c.Request = req
	// 解析form
	c.Request.ParseForm()
	// 初始化字段
	c.init(e)
	// 重置chain
	c.HandlersChain = HandlersChain{}
	e.HandleRequest(c)
	e.pool.Put(c)
}

// 启动服务
func (e *Engine) Start(addr string) (error) {
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