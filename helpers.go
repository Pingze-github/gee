package gee

import (
	"strconv"
	"strings"
)

// hostname:port 解析为 hostname和port
// 如果传入域名，则返回原hostname，port=0
func sepHostNameWithPort(host string) (hostname string, port int) {
	var err error
	index := strings.Index(hostname, ":")
	if index >= 0 {
		hostname = host[:index]
		port, err = strconv.Atoi(host[:index])
		if err != nil {
			panic(err)
		}
		return
	}
	hostname = host
	return
}
