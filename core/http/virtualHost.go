package http

import (
	"os"

	"github.com/joyme123/cats/config"
)

type VirtualHost struct {
	head    Component
	context VhostContext // VirtualHost的上下文环境
}

func (vh *VirtualHost) Init() {

	vh.context = VhostContext{make(map[string]interface{})}
}

func (vh *VirtualHost) GetContext() *VhostContext {
	return &vh.context
}

// 向VirtualHost中注入组件
func (vh *VirtualHost) Register(component Component) {
	vh.head = component
}

// Start 虚拟主机启动后调用的函数
func (vh *VirtualHost) Start() {
	vh.head.Start()
}

// 接过请求和响应的控制权
func (vh *VirtualHost) ServeHttp(req *Request, resp *Response) {

	// 解析完毕一个请求
	if config.GetInstance().Log {
		req.logger(os.Stdout)
	}

	resp.Version = req.Version

	vh.head.Serve(req, resp)

	resp.out()

	// 清空请求和响应的状态
	vh.clear(req, resp)
}

func (vh *VirtualHost) clear(req *Request, resp *Response) {
	req.Clear()
	resp.Clear()
}
