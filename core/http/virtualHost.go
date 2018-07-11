package http

import (
	"os"

	"github.com/joyme123/cats/config"
)

type VirtualHost struct {
	Request  *Request
	Response *Response
	hub      Hub
	context  Context // VirtualHost的上下文环境
}

func (vh *VirtualHost) Init() {

	vh.context = Context{make(map[string]interface{})}
}

func (vh *VirtualHost) GetContext() *Context {
	return &vh.context
}

// 向VirtualHost中注入组件
func (vh *VirtualHost) Register(component interface{}) {
	vh.hub.Register(component.(Component))
}

// 接过请求和响应的控制权
func (vh *VirtualHost) ServeHttp(req *Request, resp *Response) {
	vh.Request = req
	vh.Response = resp

	// 解析完毕一个请求
	if config.GetInstance().Log {
		vh.Request.logger(os.Stdout)
	}

	vh.Response.Version = vh.Request.Version

	for _, comp := range vh.hub.container {
		comp.Serve(vh.Request, vh.Response)
	}

	vh.Response.out()

	// 清空请求和响应的状态
	vh.clear()
}

func (vh *VirtualHost) clear() {
	vh.Request.Clear()
	vh.Response.Clear()
}
