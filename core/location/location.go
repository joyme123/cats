// location 组件可以借鉴nginx的设计，支持精确匹配，前缀匹配，正则匹配等待
package location

import (
	"log"
	"regexp"
	"strings"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

// Location 组件
type Location struct {
	Index             int
	Context           *http.VhostContext // vhost的context
	exactMatch        []Pattern          // 通用匹配
	firstPrefixMatch  []Pattern          // 在正则匹配前的前缀匹配
	regexMatch        []Pattern          // 正则匹配
	secondPrefixMatch []Pattern          // 在正则匹配后的前缀匹配
	defaultMatch      Pattern            // 通用匹配
}

// New 方法是FastCGI 的实例化
func (loc *Location) New(site *config.Site, context *http.VhostContext) {
	loc.Context = context

	loc.exactMatch = make([]Pattern, 0, 0)
	loc.firstPrefixMatch = make([]Pattern, 0, 0)
	loc.regexMatch = make([]Pattern, 0, 0)
	loc.secondPrefixMatch = make([]Pattern, 0, 0)
}

// Start 方法是FastCGI在服务启动时调用的方法
func (loc *Location) Start() {
	// 调用所有pattern中的Start

	for _, pattern := range loc.exactMatch {
		pattern.comp.Start()
	}

	for _, pattern := range loc.firstPrefixMatch {
		pattern.comp.Start()
	}

	for _, pattern := range loc.regexMatch {
		pattern.comp.Start()
	}

	for _, pattern := range loc.secondPrefixMatch {
		pattern.comp.Start()
	}

	loc.defaultMatch.comp.Start()

}

// Serve 方法是FastCGI在有请求到来时被调用的方法
func (loc *Location) Serve(req *http.Request, resp *http.Response) {

	log.Println("location serve")
	// 依次进行location的匹配
	for _, pattern := range loc.exactMatch {
		if pattern.regex == req.URI {
			// 匹配到
			pattern.comp.Serve(req, resp)
			return
		}
	}

	for _, pattern := range loc.firstPrefixMatch {
		if strings.HasPrefix(req.URI, pattern.regex) {
			// 匹配到
			pattern.comp.Serve(req, resp)
			return
		}
	}

	for _, pattern := range loc.regexMatch {

		var reg *regexp.Regexp
		if pattern.mode == 3 { // 区分大小写
			reg = regexp.MustCompile("(?i)" + pattern.regex)
		} else { // 默认就是不区分
			reg = regexp.MustCompile(pattern.regex)
		}

		if reg.MatchString(req.URI) {
			// 匹配到
			pattern.comp.Serve(req, resp)
			return
		}
	}

	for _, pattern := range loc.secondPrefixMatch {
		if strings.HasPrefix(req.URI, pattern.regex) {
			// 匹配到
			pattern.comp.Serve(req, resp)
			return
		}
	}

	// 匹配到
	loc.defaultMatch.comp.Serve(req, resp)
}

// Shutdown 方法是FastCGI在服务终止时被调用的方法
func (loc *Location) Shutdown() {

}

// GetIndex 用来获取当前组件的索引
func (loc *Location) GetIndex() int {
	return loc.Index
}

// Register 向location组件中注入其他组件
func (loc *Location) Register(modeStr string, regex string, comp interface{}) {

	var pattern Pattern
	pattern.New(modeStr, regex, comp)

	switch pattern.mode {
	case 1:
		loc.exactMatch = append(loc.exactMatch, pattern)
		break

	case 2:
		loc.firstPrefixMatch = append(loc.firstPrefixMatch, pattern)
		break

	case 3:
		loc.regexMatch = append(loc.regexMatch, pattern)
		break

	case 4:
		loc.regexMatch = append(loc.regexMatch, pattern)
		break

	case 5:
		loc.secondPrefixMatch = append(loc.secondPrefixMatch, pattern)
		break

	case 6:
		loc.defaultMatch = pattern
		break
	}
}

// GetContainer 获取父容器
func (loc *Location) GetContainer() string {
	return "vhost"
}
