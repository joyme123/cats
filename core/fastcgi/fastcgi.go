package fastcgi

import (
	"log"
	"net"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

// FastCGI 的结构体
type FastCGI struct {
	Index    int
	sockAdrr string        // fastcgi应用程序的socket地址
	Context  *http.Context // handler的context
	req      *http.Request
	resp     *http.Response
	fcgiConn net.Conn // fcgi的连接
}

// New 方法是FastCGI 的实例化
func (fcgi *FastCGI) New(site *config.Site, context *http.Context) {
	fcgi.sockAdrr = site.FCGIPass
	fcgi.Context = context
}

// Start 方法是FastCGI在服务启动时调用的方法
func (fcgi *FastCGI) Start() {

}

// Serve 方法是FastCGI在有请求到来时被调用的方法
func (fcgi *FastCGI) Serve(req *http.Request, resp *http.Response) {
	fcgi.req = req
	fcgi.resp = resp

	// 1.创建并发送开始请求
	var beginRequestRecord FCGIBeginRequestRecord
	beginRequestRecord.New(1, FCGIResponder)
	fcgi.sendRecord(&beginRequestRecord)

	// 2.获取当前请求的请求头，将其传递给fastcgi 程序
	var paramsRecord FCGIParamsRecord
	paramsRecord.New(1, fcgi.req.Headers)
	fcgi.sendRecord(&paramsRecord)

	// 3.创建并发送stdin请求
	var stdinRecord FCGIStdioRecord
	stdinRecord.New(1, []byte("a=hello&b=world"))
	fcgi.sendRecord(&stdinRecord)

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
