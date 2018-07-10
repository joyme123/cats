一些文档：
http://blog.51reboot.com/%E7%BD%91%E7%BB%9C%E7%BC%96%E7%A8%8B%EF%BC%88%E5%85%AD%EF%BC%89%EF%BC%9A%E7%AB%AF%E5%8F%A3%E9%82%A3%E4%BA%9B%E4%BA%8B%E5%84%BF/

vhost的类型

1.监听不同的ip
这种情况下，直接bind,获取新的listenfd即可。

2.监听不同的端口
这种情况下，直接获取bind,获取新的listenfd即可

3.监听相同的ip和端口
这种情况下，就需要通过host来进行转发了