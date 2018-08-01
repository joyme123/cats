package fastcgi

import (
	"log"
	"net"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/utils"
)

// FastCGI 的结构体
type FastCGI struct {
	Index    int
	sockAdrr string        // fastcgi应用程序的socket地址
	Context  *http.Context // handler的context
	req      *http.Request
	resp     *http.Response
	fcgiConn net.Conn // fcgi的连接
	RootDir  string   // 根目录地址
}

// New 方法是FastCGI 的实例化
func (fcgi *FastCGI) New(site *config.Site, context *http.Context) {
	fcgi.sockAdrr = site.FCGIPass
	fcgi.Context = context
	fcgi.RootDir = site.Root
}

// Start 方法是FastCGI在服务启动时调用的方法
func (fcgi *FastCGI) Start() {

}

// Serve 方法是FastCGI在有请求到来时被调用的方法
func (fcgi *FastCGI) Serve(req *http.Request, resp *http.Response) {

	log.Println("fastcgi serve")

	fcgi.req = req
	fcgi.resp = resp

	// 1.创建并发送开始请求
	var beginRequestRecord FCGIBeginRequestRecord
	beginRequestRecord.New(1, FCGIResponder)
	fcgi.sendRecord(&beginRequestRecord)

	// 2.获取当前请求的请求头，将其传递给fastcgi 程序
	var paramsRecord FCGIParamsRecord
	var params map[string]string
	params = make(map[string]string)

	var filepath string // 脚本文件的绝对路径
	if indexFiles, ok := fcgi.Context.KeyValue["IndexFiles"]; ok {
		filepath = utils.GetAbsolutePath(fcgi.RootDir, req.URI, indexFiles.([]string))
	} else {
		filepath = utils.GetAbsolutePath(fcgi.RootDir, req.URI, make([]string, 0, 0))
	}

	params["SCRIPT_FILENAME"] = filepath
	params["QUERY_STRING"] = "a=123?b=345"
	params["REQUEST_METHOD"] = req.Method
	params["CONTENT_TYPE"] = req.Headers["content-type"]
	params["CONTENT_LENGTH"] = req.Headers["content-length"]
	params["SCRIPT_NAME"] = req.URI
	params["REQUEST_URI"] = req.URI
	params["DOCUMENT_URI"] = req.URI
	params["DOCUMENT_ROOT"] = fcgi.RootDir
	params["SERVER_PROTOCOL"] = req.Version
	params["GATEWAY_INTERFACE"] = "CGI/1.1"
	params["SERVER_SOFTWARE"] = "cats"
	params["REMOTE_ADDR"] = "192.168.0.6"
	params["REMOTE_PORT"] = "27869"
	params["SERVER_ADDR"] = "127.0.0.1"
	params["SERVER_PORT"] = "8090"
	params["SERVER_NAME"] = "jiang"
	// params["HTTP_ACCEPT"] =

	paramsRecord.New(1, params)
	fcgi.sendRecord(&paramsRecord)

	// 3.创建并发送stdin请求
	var stdinRecord FCGIStdioRecord
	stdinRecord.New(1, []byte("a=hello&b=world"))
	fcgi.sendRecord(&stdinRecord)
	var emptyStdinRecord FCGIStdioRecord
	emptyStdinRecord.New(1, []byte(""))
	fcgi.sendRecord(&emptyStdinRecord)

}

// Shutdown 方法是FastCGI在服务终止时被调用的方法
func (fcgi *FastCGI) Shutdown() {
	log.Println("关闭fastcgi连接")
	fcgi.fcgiConn.Close()
}

// GetIndex 用来获取当前组件的索引
func (fcgi *FastCGI) GetIndex() int {
	return fcgi.Index
}

func (fcgi *FastCGI) GetContainer() string {
	return "location"
}

// sendRecord 发送FastCGI 记录
func (fcgi *FastCGI) sendRecord(record Record) {
	if fcgi.fcgiConn == nil {
		// 连接不存在,创建连接
		var err error
		fcgi.fcgiConn, err = net.Dial("tcp", fcgi.sockAdrr)
		if err != nil {
			log.Println("error when connect to FastCGI Application", err.Error())
			return
		}

		go fcgi.readHandler()

	}

	// 通过fcgiConn发送记录
	fcgi.fcgiConn.Write(record.ToBlob())
}

// readHandler 负责从FastCGI应用程序中读取stdout,stderr,以及EndRequestRecord
func (fcgi *FastCGI) readHandler() {

	for {
		var data [1024]byte
		n, err := fcgi.fcgiConn.Read(data[:])

		if err != nil {
			log.Println("error when read data from FastCGI Application", err.Error())
			return
		}

		if n == 0 {
			log.Println("没有读取数据")
		}

		// 打印出n的数据

		log.Printf("读取到长度为%d的数据：%s", n, string(data[0:n]))

	}
}
