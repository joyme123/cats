https://aotu.io/notes/2016/09/22/http-caching/index.html

http协议中虽然规定了使用etag作为文件的唯一标识符，但是并没有说具体如何实现。

在平时我们常常见到使用对文件内容做md5和sha散列作为校验，但是这在web服务器上是不推荐使用的。因为这种对全部内容做散列的方法很耗时，在高并发量的情况下，会严重影响机器的性能。于是需要参考目前情况下，一些实现机制。

## nginx 中如何实现 etag

参考的连接：https://github.com/billfeller/billfeller.github.io/issues/91

nginx中使用文件的last_modified_time和文件长度作为etag,类似于`ETag: "5b7f8254-488"`

```
 etag->value.len = ngx_sprintf(etag->value.data, "\"%xT-%xO\"",
                                  r->headers_out.last_modified_time,
                                  r->headers_out.content_length_n)
                      - etag->value.data;
```


## apache 中如何实现 etag

下面的内容来自网络：

 > 以‘-‘为分隔符，分为三节：

 > 第一节：文件inode的十六进制表示

 > 第二节：文件长度（以字节为单位）的十六进制表示

 > 第三节：文件的最后修改时间（UNIX时间戳）的十六进制表示

 >当文件跨文件系统移动时，文件inode会发生变化（当然可能有极低的概率不变化）

## nginx 和 apache 的实现方式有什么问题

nignx下，在多服务器负载均衡的时候，可能因为部署过程是串行的，所以上传上去的时间是不同的。导致同一个资源在不同的服务器上last_modified_time是不同的。这样生成的etag生成可能就不一致。所以如果需要使用etag，则需要一定的策略来保证各个服务器上文件的last_modified_timem是一样的。

apache下，因为apache中的实现还增加了文件的inode信息。这个inode是文件在文件系统中的信息，因此如果apache集群中使用etag，必须禁止使用inode，同时保证last_modified_time一致。

