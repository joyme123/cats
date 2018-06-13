# cats

项目结构

--config            存放配置解析

--core              核心部分代码

 - --http          http协议的实现，不准备使用go的官方库，自己从tcp部分开始。计划只支持http1.1和http2协议

 - --index         index模块的实现，支持索引文件的配置

 - --gzip          gzip模块的实现，支持压缩

 - --image         image模块的实现，支持图片的基本处理

 - --cache         cache模块的实现，支持缓存

 - --fastcgi       fastcgi协议的支持，用来支持php-fpm

 - --header        用来支持自定义头

 - --tls           用来支持https
 

# 下一步开发计划

1.vhost 支持多站点(多端口或者单端口多域名)


linux设置ulimit

    https://blog.csdn.net/bugall/article/details/45869183
