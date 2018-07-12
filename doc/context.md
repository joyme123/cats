# 关于context的设计

## 1.context概念简介

context是上下文的意思，在很多系统的设计中都会有context，在大多数系统中，context都是单例的，因为context中保存的是整个系统在运行中的一些共享变量。在cats的实现中，context不是单例的,这是因为cats是支持VirtualHost(下文简称为vhost)的，也就是说会同时支持多个虚拟主机的。因此context是一个vhost的上下文环境，对于一个vhost来说，它是单例的。


## 2.context是如何在一个vhost中传递的？

context在一个vhost初始化时被注入，在vhost中的组件被实例化时，会注入当前vhost的context指针，这样，vhost及其组件都拥有了context的引用，可以获取和修改context中的内容。

## 3.context的具体作用是什么？

Context其实是一个很简单的结构体，目前只包含一个map，主要作用是通过一个string设置和获取其对应的值。

```
    type Context struct {
        KeyValue map[string]interface{}
    }
```

具体作用有以下：

- index组件、serveFile组件：index组件会向其中写入IndexFiles，也就是当url不指定具体资源时，默认加载的文件。index组件写入后，在之后的serveFile组件会根据这个值去寻找应该读取的文件

## 4.一些要注意的点

 - 1.vhost的context中，应该只保存vhost中共享的，如果是要在一个request请求中共享一些值，应该使用request中的context，这是针对于一次请求的上下文
