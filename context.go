package gee

import (
	"net/http"
	"net/url"
)

// 上下文
// TODO 一个上下文中会有多个handler，逐个执行
type Context struct {
	isEnd bool
	ResponseWriter http.ResponseWriter
	Request *http.Request
	HandlersChain HandlersChain
	FinalHandler func(*Context, interface{})

	// Request
	Method string
	Header http.Header
	Proto string
	Addr string
	Host string
	HostName string
	Url string
	Path string
	Query url.Values
	Form url.Values
	Ip string
	Ips []string

	// Response
	Status int

	// Status
	Wrote bool
}

type HandlersChain []GeeHandler

// [pri] 从req/res中获取主要属性
func (c *Context) init(e *Engine) {
	// 统一处理err
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	var err error

	c.FinalHandler = e.FinalHandler

	c.Method = c.Request.Method
	c.Header = c.Request.Header
	c.Proto = c.Request.Proto
	c.Addr = c.Request.RemoteAddr
	c.Host = c.Request.Host
	c.HostName, _ = sepHostNameWithPort(c.Host)
	c.Url = c.Request.RequestURI
	c.Path = c.Request.URL.Path
	c.Query, err = url.ParseQuery(c.Request.URL.RawQuery);
	if err != nil {
		panic(err)
	}
	c.Request.ParseForm()
	c.Form = c.Request.PostForm
	c.Ip, _ = sepHostNameWithPort(c.Addr)
	// TODO Ips
}

// [pri] 获取下一个handler
func (c *Context) getNextHandler() (GeeHandler) {
	handler := c.HandlersChain[0]
	c.HandlersChain = c.HandlersChain[1:]
	return handler
}

// 结束检查
func (c *Context) CheckEnd() {
	//if c.isEnd {
	//	panic(errors.New("[gee] Can not write to gee.Context after end it"))
	//}
}

// 结束
func (c *Context) End(text string) (error) {
	c.CheckEnd()
	var err error
	if text != "" {
		_, err = c.ResponseWriter.Write([]byte(text))
	}
	c.isEnd = true
	return err
}

// 继续向下一个中间件执行
func (c *Context) Next() {
	if len(c.HandlersChain) == 0 {
		return
	}
	handler := c.getNextHandler()
	handler.Serve(c)
}

// 先执行context内部序列中的其他方法
// 可以实现类似koa框架的 middeware - handler - middeware 结构
func (c *Context) Yield() {
	if len(c.HandlersChain) == 0 {
		return
	}
	handler := c.getNextHandler()
	handler.Serve(c)
}

// 立即跳转到FinalHandler进行最终处理
func (c *Context) Final(data interface{}) {
	c.FinalHandler(c, data)
}

// 写数据
func (c *Context) Write(bytes []byte) (int, error) {
	defer func() {
		c.Wrote = true
		c.Status = 200
	}()
	c.CheckEnd()
	return c.ResponseWriter.Write(bytes)
}

// 写数据(string)
func (c *Context) WriteString(text string) (int, error) {
	c.CheckEnd()
	return c.ResponseWriter.Write([]byte(text))
}

// 设置状态码
func (c *Context) SetStatus(status int) {
	c.Status = status
	c.ResponseWriter.WriteHeader(404)
}

// 解析json
func (c *Context) ParseFormJson(t interface{}) (err error) {
	err = parseFormJson(c, t)
	return
}