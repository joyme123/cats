package fastcgi

import (
	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

// FastCGI 的结构体
type FastCGI struct {
	Index    int
	sockAdrr string        // fastcgi应用程序的socket地址
	Context  *http.Context // handler的context
	req      *http.Request
	resp     *http.Response
}

// New 方法是FastCGI 的实例化
func (fcgi *FastCGI) New(site *config.Site, context *http.Context) {
	fcgi.sockAdrr = site.FCgiPass
	fcgi.Context = context
}

// Start 方法是FastCGI在服务启动时调用的方法
func (fcgi *FastCGI) Start() {

}

// Serve 方法是FastCGI在有请求到来时被调用的方法
func (fcgi *FastCGI) Serve(req *http.Request, resp *http.Response) {
	fcgi.req = req
	fcgi.resp = resp

	// 获取当前请求的请求头，将其传递给fastcgi 程序

}

// Shutdown 方法是FastCGI在服务终止时被调用的方法
func (fcgi *FastCGI) Shutdown() {

}

// GetIndex 用来获取当前组件的索引
func (fcgi *FastCGI) GetIndex() int {
	return fcgi.Index
}
