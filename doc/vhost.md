## vhost的类型

1.监听不同的ip
这种情况下，直接bind,获取新的listenfd即可。

2.监听不同的端口
这种情况下，直接获取bind,获取新的listenfd即可

3.监听相同的ip和端口
这种情况下，就需要通过host来进行转发了

## 如何实现，一个端口下有多个vhost

一个端口有一个socket监听，代表着一个server，一个server下有多个vhost,server负责初始化Response和Request，然后根据Request中的Host,将Response和Request交给对应的vhost控制。server的任务到此结束。之后vhost掌管Response和Request