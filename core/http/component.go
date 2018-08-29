// http 包, 封装了plugin的注册,执行机制
package http

import (
	"github.com/joyme123/cats/config"
)

// 所有的plugin都需要实现这个接口
type Component interface {

	// 组件初始化,注入VirtualHost的上下文环境
	New(site *config.Site, context *VhostContext)

	// 在服务启动时执行
	Start()

	// 在有请求到来时执行
	Serve(req *Request, resp *Response)

	// 在服务关闭时执行
	Shutdown()

	// 设置下一个组件
	SetNext(nextComp Component)

	// 调用下一个组件
	Next(req *Request, resp *Response)

	// 获取index, index的作用是指定插件的执行顺序
	GetIndex() int

	// 获取组件的寄主身份，比如Index,Mime组件应该是属于vhost的，而serveFile和fastcgi是属于location
	GetContainer() string
}

// BaseComponent 是所有组件的基类
type BaseComponent struct {
	next    Component
	hasNext bool
}

// SetNext 用来设置该组件下一个要调用的组件
func (comp *BaseComponent) SetNext(nextComp Component) {
	comp.next = nextComp
	comp.hasNext = true
}

// Next 用来调用下一个组件, 子类可以选择重写该方法
func (comp *BaseComponent) Next(req *Request, resp *Response) {
	if comp.hasNext {
		comp.next.(Component).Serve(req, resp)
	}
}
