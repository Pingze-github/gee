package gee

import (
	"regexp"
)

type Route struct {
	Method string
	Path string
	Pattern *regexp.Regexp
}

// 创建路由
func createRoute(method string, path string) (r *Route, err error) {
	var patternString string
	var pattern *regexp.Regexp

	// 将一般表示的路径，改为正则格式
	patternString = regexp.MustCompile(`\*`).ReplaceAllString(path, ".*")

	patternString = "^" + patternString + "$"

	pattern, err = regexp.Compile(patternString)
	if err != nil {
		return &Route{}, err
	}
	r = &Route{
		Method: method,
		Path: path,
		Pattern: pattern,
	}
	return
}

// 判断路由是否符合
func (r *Route) Match(c *Context) bool {
	if r.Pattern.MatchString(c.Request.URL.Path) &&
		(r.Method == c.Request.Method || r.Method == "ALL") {
		return true
	}
	return false
}


