package config

// Site，用来实例化一个vhost实例
type Site struct {
	Addr       string
	Port       int
	ServerName string
	Root       string
	Index      []string
	FCGIPass   string            // fcgi程序的socket地址
	FCGIParams map[string]string // 要传给fcgi程序的params，可配置
}

// FCGIParams 可配置项有:

/******
以$符号开头的，是内置的变量，否则就是字符串
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
fastcgi_param  SERVER_SOFTWARE    cats;

fastcgi_param  REMOTE_ADDR        $remote_addr;
fastcgi_param  REMOTE_PORT        $remote_port;
fastcgi_param  SERVER_ADDR        $server_addr;
fastcgi_param  SERVER_PORT        $server_port;
fastcgi_param  SERVER_NAME        $server_name;

fastcgi_param  REDIRECT_STATUS    200;

fastcgi_param  PATH_INFO		  ^(.+\.php)(/.+)$

*/
