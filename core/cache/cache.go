package cache

import (
	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type Cache struct {
	Index int
}

// 组件初始化,注入VirtualHost的上下文环境
func (cache *Cache) New(site *config.Site, context http.VhostContext) {

}

// 在服务启动时执行
func (cache *Cache) Start() {

}

// Serve 获取req中的 If-Modified-Since,If-Unmodified-Since,If-None-Match,If-Match
func (cache *Cache) Serve(req *http.Request, resp *http.Response) {

}

// 在服务关闭时执行
func (cache *Cache) Shutdown() {

}

// 获取index, index的作用是指定插件的执行顺序
func (cache *Cache) GetIndex() int {
	return cache.Index
}

// 获取组件的寄主身份，比如Index,Mime组件应该是属于vhost的，而serveFile和fastcgi是属于location
func (cache *Cache) GetContainer() string {
	return "location"
}
