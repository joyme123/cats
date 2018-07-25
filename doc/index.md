# 关于cats服务器

cats服务器是用go语言实现的一个http服务器。没有使用go语言的http库，是从tcp层开始实现的。

![cats服务的结构图](https://github.com/joyme123/cats/raw/master/doc/img/cats-structure.jpg)

在cats被启动后，会根据配置文件中的端口号进行监听，比如在上图中有3个端口号需要监听：80,8080,8090。每个端口号对应着一个server，每个server下都有多个VirtualHost，也就是虚拟主机。每个虚拟主机是通过server_name区分开的。

80端口对应着server1，server1下有两个虚拟主机，VirtualHost1是默认的服务，也就是说，我们通过ip访问80端口，会默认访问到VirtualHost1下的Web服务。如果我们访问myway5.com，也会访问到VirtualHost1，如果我们访问sub.myway5.com，则会访问到VirtualHost2。

在机器性能和端口号足够的条件下，可以有尽可能多的server，一个server下也可以有尽可能多的VirtualHost。

## 下面会记录一些在开发过程中的思考，设计等等细节

 - [vhost的设计](vhost.md)

 - [context的设计](context.md)
 
 - [component的设计](component_design.md)
 
 - [一个标准的http请求的处理过程]()