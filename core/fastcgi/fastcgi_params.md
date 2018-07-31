


fastcgi程序在和web服务器交互时，需要传递很多的参数。

http://php.net/manual/zh/reserved.variables.server.php

这是一份来自网络的参数列表,后面标注1的都是和nginx的参数列表重复的。经过抓包查看，这一份在和php交互时是最完整的

SCRIPT_FILENAME,            // 脚本的绝对地址 /home/jiang/php/index.php

QUERY_STRING, 1             // 请求的query_string,a=b&c=d

REQUEST_METHOD,1            // 请求的request_method

CONTENT_TYPE,1              // 请求的content_type

CONTENT_LENGTH,1            // 请求的content_length

SCRIPT_NAME,1               // 脚本名, /php/index.php

REQUEST_URI,1               // 请求地址,比如/php/index.php?a=b&c=d

DOCUMENT_URI,1              // 请求的文件地址，也就是除去域名之后的地址比如/php/index.php

DOCUMENT_ROOT,1             // 这里就是server的Root参数

SERVER_PROTOCOL,1           // HTTP协议版本, Http/1.1

GATEWAY_INTERFACE,1         // CGI协议版本，CGI/1.1

SERVER_SOFTWARE,1           // Web服务器的名称

REMOTE_ADDR,1               // http请求的ip

REMOTE_PORT,1               // 用户机器上连接到 Web 服务器所使用的端口号。

SERVER_ADDR,1               // 当前运行脚本所在的服务器的 IP 地址。

SERVER_PORT,1               // Web 服务器使用的端口。

SERVER_NAME,1               // 当前运行脚本所在的服务器的主机名

HTTP_ACCEPT,1

HTTP_ACCEPT_LANGUAGE,

HTTP_ACCEPT_ENCODING,

HTTP_USER_AGENT,

HTTP_HOST,

HTTP_CONNECTION,

HTTP_CONTENT_TYPE,

HTTP_CONTENT_LENGTH,

HTTP_CACHE_CONTROL,

HTTP_COOKIE,

HTTP_FCGI_PARAMS_MAX        // 这个在实际实验中没有看到

还有个额外的PATH_INFO, PATH_INFO 是在类似这种情况下：http://a.com/index.php/a/b/c/d 那么后面的/a/b/c/d就是path_info

https://blog.jjonline.cn/linux/218.html

下面是PHP独有的
REDIRECT_STATUS,1    Http请求的状态置


下面是一些来自nginx配置文件中的fastcgi_params

fastcgi_param  QUERY_STRING       $query_string;
fastcgi_param  REQUEST_METHOD     $request_method;
fastcgi_param  CONTENT_TYPE       $content_type;
fastcgi_param  CONTENT_LENGTH     $content_length;

fastcgi_param  SCRIPT_NAME        $fastcgi_script_name;
fastcgi_param  REQUEST_URI        $request_uri;
fastcgi_param  DOCUMENT_URI       $document_uri;
fastcgi_param  DOCUMENT_ROOT      $document_root;
fastcgi_param  SERVER_PROTOCOL    $server_protocol;
fastcgi_param  REQUEST_SCHEME     $scheme;
fastcgi_param  HTTPS              $https if_not_empty;

fastcgi_param  GATEWAY_INTERFACE  CGI/1.1;
fastcgi_param  SERVER_SOFTWARE    nginx/$nginx_version;

fastcgi_param  REMOTE_ADDR        $remote_addr;
fastcgi_param  REMOTE_PORT        $remote_port;
fastcgi_param  SERVER_ADDR        $server_addr;
fastcgi_param  SERVER_PORT        $server_port;
fastcgi_param  SERVER_NAME        $server_name;

# PHP only, required if PHP was built with --enable-force-cgi-redirect
fastcgi_param  REDIRECT_STATUS    200;