package config

var context *Context

type Context struct {
	Log bool
}

func GetInstance() *Context {
	if context == nil {
		context = &Context{Log: true}
	}

	return context
}
