# cats

## 项目简介

cats 是一个用go语言写的web server。没有使用go官方的http库，http方面的代码都是从tcp层开始写的。开发cats的目的只是出于个人的兴趣，cats的开发方向是以**简单易用**为原则的，在此基础上会做一些有趣的功能。

项目的整体架构如下图

![cats服务的结构图](https://github.com/joyme123/cats/raw/master/doc/img/cats-structure.jpg)

相关说明见文档：[cats的相关文档](https://github.com/joyme123/cats/blob/master/doc/index.md)

## 项目编译和运行

```
# 编译
go build main.go

# 运行
go run main.go
```

## 项目文档

这个项目中有一些开发时的笔记，存放在`doc`目录下。也可以直接访问链接查看:[cats的相关文档](https://github.com/joyme123/cats/blob/master/doc/index.md)。

希望这些文档可以对你有帮助

## 项目结构

--config            存放配置解析

--core              核心部分代码

 - --http          http协议的实现，不准备使用go的官方库，自己从tcp部分开始。计划只支持http1.1和http2协议

 - --index         index模块的实现，支持索引文件的配置

 - --gzip          gzip模块的实现，支持压缩. TODO

 - --image         image模块的实现，支持图片的基本处理 TODO

 - --cache         cache模块的实现，支持缓存 In progress

 - --fastcgi       fastcgi协议的支持，用来支持php-fpm

 - --header        用来支持自定义头 TODO

 - --tls           用来支持https TODO
 
 - --serveFile     用来serve静态文件

 - --mime          用来生成mime

 - --location      根据url进行location
 
-- utils           存放了一些通用的工具类

## LICENSE

MIT