//Package http 是对http协议的封装
package http

import (
	"github.com/joyme123/cats/config"
)

//Server 类，目前只支持http1.1和http2
//addr是服务的监听地址
//port是监听的端口
type Server struct {
	srvname string
	configs map[string]vhost
}

//vhost,是虚拟Host的实现
type vhost struct {
	config config.Config
}
