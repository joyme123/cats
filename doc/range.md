https://tools.ietf.org/html/rfc7233#section-3.2

The Range 是一个请求首部，告知服务器返回文件的哪一部分。在一个  Range 首部中，可以一次性请求多个部分，服务器会以 multipart 文件的形式将其返回。如果服务器返回的是范围响应，需要使用 206 Partial Content 状态码。假如所请求的范围不合法，那么服务器会返回  416 Range Not Satisfiable 状态码，表示客户端错误。服务器允许忽略  Range  首部，从而返回整个文件，状态码用 200 。


```
Range: <unit>=<range-start>-
Range: <unit>=<range-start>-<range-end>
Range: <unit>=<range-start>-<range-end>, <range-start>-<range-end>
Range: <unit>=<range-start>-<range-end>, <range-start>-<range-end>, <range-start>-<range-end>

```

## 指令

<unit>
范围所采用的单位，通常是字节（bytes）。
<range-start>
一个整数，表示在特定单位下，范围的起始值。
<range-end>
一个整数，表示在特定单位下，范围的结束值。这个值是可选的，如果不存在，表示此范围一直延伸到文档结束。

## 示例

Range: bytes=200-1000, 2000-6576, 19000-

## Range 请求的响应
https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Range_requests

```
HTTP/1.1 206 Partial Content
Content-Type: multipart/byteranges; boundary=3d6b6a416f9b5
Content-Length: 282

--3d6b6a416f9b5
Content-Type: text/html
Content-Range: bytes 0-50/1270

<!doctype html>
<html>
<head>
    <title>Example Do
--3d6b6a416f9b5
Content-Type: text/html
Content-Range: bytes 100-150/1270

eta http-equiv="Content-type" content="text/html; c
--3d6b6a416f9b5--
```

## 文档地址

https://tools.ietf.org/html/rfc7233#section-3.1


使用curl测试range请求
curl -H  "Range: bytes=0-50, 100-150" -X GET http://mysite.com:8089/php/test_fcgi_request.html
