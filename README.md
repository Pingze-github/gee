# gee
My HTTP web framework written in GO.

## TODO
1. ~~将(rw http.ResponseWriter, req *http.Request) 转化为 (c *gee.Context)~~
2. ~~支持直接注册func而非struct~~
3. ~~支持绑定不同ip~~
4. 支持httprouter高级路由
5. ~~支持中间件~~
6. ~~支持最终(异常)处理中间件~~
7. ~~性能测试，探究并发性能。不同context间的并发性在net/http中实现~~
8. ~~在已被占用的端口启动时，异常退出~~
9. Render为JSON字符串的实现
10. ~~Context基本属性和方法~~
11. ~~请求参数的解析~~
12. ~~支持同时传入多个handlers~~
13. ~~支持content-type: application/json~~