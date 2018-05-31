// hub包, 封装了plugin的注册,执行机制
package hub

import (
	"github.com/joyme123/cats/core/http"
)

// 所有的plugin都需要实现这个接口
type Controller interface {

	// 在服务启动时执行
	Start()

	// 在有请求到来时执行
	Serve(req *http.Request, resp *http.Response)

	// 在服务关闭时执行
	Shutdown()

	// 获取index, index的作用是指定插件的执行顺序
	getIndex()
}
