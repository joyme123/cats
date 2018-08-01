// location 组件的匹配规则
package location

import "github.com/joyme123/cats/core/http"

// Pattern 是 location组件的匹配模式
// mode代表当前是什么匹配，1是精确匹配，2是前缀匹配，3是按文件中顺序的正则匹配，4匹配不带任何修饰的前缀匹配,5 通用匹配
// https://moonbingbing.gitbooks.io/openresty-best-practices/ngx/nginx_local_pcre.html

type Pattern struct {
	mode  int
	regex string

	// FIXME: 这个组件应该是指针
	comp http.Component
}

// New 用来实例化一个Pattern对象
func (pattern *Pattern) New(modeStr string, regex string, comp interface{}) {
	switch modeStr {
	case "=": //精确匹配
		pattern.mode = 1
		break

	case "^~": // 在正则前的前缀匹配
		pattern.mode = 2
		break

	case "~": // 区分大小写的正则匹配
		pattern.mode = 3
		break

	case "~*": // 不区分大小写的正则匹配
		pattern.mode = 4
		break

	case "": // 在正则后的前缀匹配
		pattern.mode = 5
		break

	case "/": // 通用匹配
		pattern.mode = 6
		break
	}

	pattern.regex = regex
	pattern.comp = comp.(http.Component)
}
