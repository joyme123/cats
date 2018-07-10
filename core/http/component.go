// http 包, 封装了plugin的注册,执行机制
package http

import "github.com/joyme123/cats/config"

// 所有的plugin都需要实现这个接口
type Component interface {

	// 组件初始化
	New(context *Context, vhost *config.VHost)

	// 在服务启动时执行
	Start()

	// 在有请求到来时执行
	Serve(req *Request, resp *Response)

	// 在服务关闭时执行
	Shutdown()

	// 获取index, index的作用是指定插件的执行顺序
	GetIndex() int
}
