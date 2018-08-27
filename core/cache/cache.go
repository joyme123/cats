package cache

import (
	"log"
	"strings"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type Cache struct {
	Index   int
	Context *http.VhostContext // handler的context
}

// 组件初始化,注入VirtualHost的上下文环境
func (cache *Cache) New(site *config.Site, context *http.VhostContext) {
	cache.Context = context
}

// 在服务启动时执行
func (cache *Cache) Start() {

}

// Serve 获取req中的 If-Modified-Since,If-Unmodified-Since,If-None-Match,If-Match
// etag 优先级要大于 last-modified
func (cache *Cache) Serve(req *http.Request, resp *http.Response) {

	filepath, ok := req.Context["FilePath"]

	if !ok {
		log.Println("error when read filepath")
	}

	// GET 和 HEAD 中，判断if-none-match和If-Modified-Since
	// RFC7232：https://tools.ietf.org/html/rfc7232#section-3.1
	if req.Method == "GET" || req.Method == "HEAD" {
		if reqEtag, ok := req.Headers["if-none-match"]; ok {
			// etag有三种方式，单值,列表,*
			if reqEtag == "*" {
				resp.AppendHeader("etag", reqEtag)
				// 返回304
				resp.Status304()
				return
			}

			// 根据etag查看更新
			etagStr := Etag(filepath.(string))

			etagList := strings.Split(etagStr, ",")
			for _, etag := range etagList {
				if etag == reqEtag {
					resp.AppendHeader("etag", reqEtag)
					// 返回304
					resp.Status304()
				}
			}

		} else if reqDate, ok := req.Headers["if-modified-since"]; ok {
			if !CompareFileModifiedTime(filepath.(string), reqDate) {
				// 早于或等于
				resp.AppendHeader("etag", reqDate)
				// 返回304
				resp.Status304()
				return
			}
		}
	}

	// TODO: if-match和if-unmodified-since需要单独处理
	// if reqEtag, ok := req.Headers["if-match"]; ok {
	// 	// etag有三种方式，单值,列表,*
	// 	if reqEtag == "*" {

	// 	} else {
	// 		// 根据etag查看更新
	// 		etag := Etag(filepath.(string))
	// 	}
	// } else if reqDate, ok := req.Headers["if-unmodified-since"]; ok {

	// 	// reqDate只有一个GMT时间戳

	// } else {
	// 	// 如果都没有的话，为当前静态文件生成etag和last-modified
	// }

	// TODO:这里还要考虑if-range，但是目前还没有实现分段请求的功能
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
	return "vhost"
}
