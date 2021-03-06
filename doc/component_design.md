## 组件（component）的设计

一个虚拟主机（VirtualHost）包含了各种各样的组件，每个组件在VirtualHost初始化的时候被实例化。之后在每次Http请求到达的时候，都会被调用。为了完成从实例化到请求时调用的通用性，每个组件都需要有一些通用的对外接口。具体的接口定义如下(v0.0.1版本)：

```
// 所有的plugin都需要实现这个接口
type Component interface {

	// 组件初始化,注入VirtualHost的上下文环境
	New(site *config.Site, context *Context)

    // 在服务启动时执行
	Start()

	// 在有请求到来时执行
	Serve(req *Request, resp *Response)

	// 在服务关闭时执行
	Shutdown()

	// 获取index, index的作用是指定组件的执行顺序
	GetIndex() int
}
```

这样，在每个组件被初始化时，系统会调用New方法将配置文件和当前虚拟主机的上下文环境注入到组件中。

在虚拟主机启动完毕的时候，每个组件的Start都会被执行一次（在整个虚拟主机的生命周期中，只会调用一次）。Start方法中，主要用来做一些初始化操作，比如FastCGI组件，可以建立起与FastCGI应用程序的连接。

在虚拟主机即将关闭之前，每个组件的Shutdown都会被执行一次（在整个虚拟主机的生命周期中，只会调用一次）。Shutdown方法中，主要用来做一些释放资源的操作，比如FastCGI组件，可以断开和FastCGI应用程序的连接。

至于Serve方法，则是在每次请求到来的时候被调用，它会被注入Request和Response的指针，用来读取Request和修改Response的内容。比如在mime组件中，mini组件会读取本次Request的文件的后缀，向Response中写入对应的content-type头。

GetIndex目前是用来指定组件的执行顺序。通过GetIndex为每个组件排好序，在每一次请求到来时，会按照这个顺序来依次启动组件。