package http

import (
	"os"

	"github.com/joyme123/cats/config"
)

type VirtualHost struct {
	hub     Hub
	context VhostContext // VirtualHost的上下文环境
}

func (vh *VirtualHost) Init() {

	vh.context = VhostContext{make(map[string]interface{})}
}

func (vh *VirtualHost) GetContext() *VhostContext {
	return &vh.context
}

// 向VirtualHost中注入组件
func (vh *VirtualHost) Register(component interface{}) {
	vh.hub.Register(component.(Component))
}

// 接过请求和响应的控制权
func (vh *VirtualHost) ServeHttp(req *Request, resp *Response) {

	// 解析完毕一个请求
	if config.GetInstance().Log {
		req.logger(os.Stdout)
	}

	resp.Version = req.Version

	for _, comp := range vh.hub.container {
		comp.Serve(req, resp)
	}

	resp.out()

	// 清空请求和响应的状态
	vh.clear(req, resp)
}

func (vh *VirtualHost) clear(req *Request, resp *Response) {
	req.Clear()
	resp.Clear()
}
