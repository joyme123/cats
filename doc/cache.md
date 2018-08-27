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


## If-Match（来自MDN）

请求首部 If-Match 的使用表示这是一个条件请求。在请求方法为 GET 和 HEAD 的情况下，服务器仅在请求的资源满足此首部列出的 ETag 之一时才会返回资源。而对于 PUT 或其他非安全方法来说，只有在满足条件的情况下才可以将资源上传。

The comparison with the stored ETag 之间的比较使用的是强比较算法，即只有在每一个比特都相同的情况下，才可以认为两个文件是相同的。在 ETag 前面添加    W/ 前缀表示可以采用相对宽松的算法。

以下是两个常见的应用场景：

 - For GET  和 HEAD 方法，搭配  Range首部使用，可以用来保证新请求的范围与之前请求的范围是对同一份资源的请求。如果  ETag 无法匹配，那么需要返回 416 (Range Not Satisfiable，范围请求无法满足) 响应。

 - 对于其他方法来说，尤其是 PUT, If-Match 首部可以用来避免更新丢失问题。它可以用来检测用户想要上传的不会覆盖获取原始资源之后做出的更新。如果请求的条件不满足，那么需要返回  412 (Precondition Failed，先决条件失败) 响应。

> Header type	Request header

> Forbidden header name	no

### 语法

```
If-Match: <etag_value>
If-Match: <etag_value>, <etag_value>, …
If-Match: *
```

### 示例

```
If-Match: "bfc13a64729c4290ef5b2c2730249c88ca92d82d"

If-Match: W/"67ab43", "54ed21", "7892dd"

If-Match: *
```

## If-None-Match（来自MDN）

If-None-Match 是一个条件式请求首部。对于 GETGET 和 HEAD 请求方法来说，当且仅当服务器上没有任何资源的 ETag 属性值与这个首部中列出的相匹配的时候，服务器端会才返回所请求的资源，响应码为  200  。对于其他方法来说，当且仅当最终确认没有已存在的资源的  ETag 属性值与这个首部中所列出的相匹配的时候，才会对请求进行相应的处理。

对于  GET 和 HEAD 方法来说，当验证失败的时候，服务器端必须返回响应码 304 （Not Modified，未改变）。对于能够引发服务器状态改变的方法，则返回 412 （Precondition Failed，前置条件失败）。需要注意的是，服务器端在生成状态码为 304 的响应的时候，必须同时生成以下会存在于对应的 200 响应中的首部：Cache-Control、Content-Location、Date、ETag、Expires 和 Vary 。

ETag 属性之间的比较采用的是弱比较算法，即两个文件除了每个比特都相同外，内容一致也可以认为是相同的。例如，如果两个页面仅仅在页脚的生成时间有所不同，就可以认为二者是相同的。

当与  If-Modified-Since  一同使用的时候，If-None-Match 优先级更高（假如服务器支持的话）。

以下是两个常见的应用场景：

 - 采用 GET 或 HEAD  方法，来更新拥有特定的ETag 属性值的缓存。
 - 采用其他方法，尤其是  PUT，将 If-None-Match used 的值设置为 * ，用来生成事先并不知道是否存在的文件，可以确保先前并没有进行过类似的上传操作，防止之前操作数据的丢失。这个问题属于更新丢失问题的一种。

> Header type	Request header

> Forbidden header name	no

### 语法

```
If-None-Match: <etag_value>
If-None-Match: <etag_value>, <etag_value>, …
If-None-Match: *
```

### 示例

```
If-None-Match: "bfc13a64729c4290ef5b2c2730249c88ca92d82d"

If-None-Match: W/"67ab43", "54ed21", "7892dd"

If-None-Match: *
```