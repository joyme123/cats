package fastcgi

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/utils"
)

// FastCGI 的结构体
type FastCGI struct {
	Index      int
	sockAdrr   string        // fastcgi应用程序的socket地址
	Context    *http.Context // handler的context
	req        *http.Request
	resp       *http.Response
	fcgiConn   net.Conn // fcgi的连接
	RootDir    string   // 根目录地址
	RequestID  uint16   // 请求id，这里现在因为只有一个连接，所以使用自增的做法
	finishChan chan bool
}

func (fcgi *FastCGI) commonHeaders() {
	fcgi.resp.AppendHeader("connection", "keep-alive")
	fcgi.resp.AppendHeader("server", "cats")
}

// New 方法是FastCGI 的实例化
func (fcgi *FastCGI) New(site *config.Site, context *http.Context) {
	fcgi.sockAdrr = site.FCGIPass
	fcgi.Context = context
	fcgi.RootDir = site.Root
	fcgi.RequestID = 1
	fcgi.finishChan = make(chan bool)
}

// Start 方法是FastCGI在服务启动时调用的方法
func (fcgi *FastCGI) Start() {

}

// Serve 方法是FastCGI在有请求到来时被调用的方法
func (fcgi *FastCGI) Serve(req *http.Request, resp *http.Response) {

	log.Println("fastcgi serve")

	fcgi.req = req
	fcgi.resp = resp

	// 建立连接
	if fcgi.fcgiConn == nil {
		// 连接不存在,创建连接
		var err error
		fcgi.fcgiConn, err = net.Dial("tcp", fcgi.sockAdrr)
		if err != nil {
			log.Println("error when connect to FastCGI Application", err.Error())
			fcgi.fcgiConn = nil
			return
		}

		log.Println("新建的fcgi连接")

		go fcgi.readHandler()
	}

	currentID := fcgi.RequestID
	log.Printf("请求ID:%v\n", currentID)
	fcgi.RequestID++
	// 1.创建并发送开始请求
	var beginRequestRecord FCGIBeginRequestRecord
	beginRequestRecord.New(currentID, FCGIResponder)

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
	params["QUERY_STRING"] = req.QueryString
	params["REQUEST_METHOD"] = req.Method
	params["CONTENT_TYPE"] = req.Headers["content-type"]

	if contentLength, ok := req.Headers["content-length"]; !ok {
		params["CONTENT_LENGTH"] = "0"
	} else {
		params["CONTENT_LENGTH"] = contentLength
	}
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

	paramsRecord.New(currentID, params)
	fcgi.sendRecord(&paramsRecord)

	var emptyParamsRecord FCGIParamsRecord
	emptyParams := make(map[string]string)
	emptyParamsRecord.New(currentID, emptyParams)
	fcgi.sendRecord(&emptyParamsRecord)

	// 3.创建并发送stdin请求
	var stdinRecord FCGIStdinRecord
	log.Printf("请求体的内容:%s", req.Body)
	stdinRecord.New(currentID, req.Body)
	fcgi.sendRecord(&stdinRecord)

	if len(req.Body) != 0 {
		var emptyStdinRecord FCGIStdinRecord
		emptyBytes := make([]byte, 0)
		emptyStdinRecord.New(currentID, emptyBytes)
		fcgi.sendRecord(&emptyStdinRecord)
	}

	// 这里应该阻塞起来等待fastcgi程序响应

	<-fcgi.finishChan // 使用管道阻塞

	log.Println("阻塞结束")
	fcgi.commonHeaders()
	fcgi.resp.StatusCode = 200
	fcgi.resp.Desc = "OK"
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
	// 通过fcgiConn发送记录
	_, err := fcgi.fcgiConn.Write(record.ToBlob())

	if err != nil {
		log.Printf("fcgi 写入错误: %v\n", err)
	}

}

// readHandler 负责从FastCGI应用程序中读取stdout,stderr,以及EndRequestRecord
func (fcgi *FastCGI) readHandler() {
	var header FCGIHeader
	readLen := 8 // 下一次要读取的长度
	isHeader := true
	stdoutHeader := false // 是否解析过标准输出的头信息
	for {
		data := make([]byte, readLen, readLen)
		n, err := fcgi.fcgiConn.Read(data[:])

		if err != nil {
			fcgi.fcgiConn.Close() // 读取出错，可能是对方已经关闭了写通道,直接关闭连接
			fcgi.fcgiConn = nil   // 置为空
			fcgi.RequestID = 1
			log.Println("error when read data from FastCGI Application", err.Error())
			return
		}

		if n == 0 {
			log.Println("没有读取数据")
		}

		// 打印出n的数据

		log.Printf("读取到长度为%d的数据：%s", n, string(data[0:n]))

		if isHeader {
			header.New(data[0:readLen]) // 初始化header
			isHeader = false
			readLen = int(header.ContentLength) + int(header.PaddingLength)
		} else {
			switch header.Type {
			case FCGIStdout: // 标准输出流
				outdata := data[0:readLen]

				// 将outdata解析出来
				state := 0 // 0代表处于普通字符状态，1代表处于\r\n状态
				var pre byte
				start := 0
				isFinish := false

				if stdoutHeader {
					if fcgi.resp.Body != nil {
						fcgi.resp.Body = append(fcgi.resp.Body, outdata[:header.ContentLength]...)
					} else {
						fcgi.resp.Body = outdata[:header.ContentLength]
					}
				} else {
					for i, v := range outdata {
						switch v {
						case '\r':
							break
						case '\n':
							if pre == '\r' {
								if state == 0 {
									// 还在解析头部
									headerStr := string(outdata[start : i-1])
									resStr := strings.Split(headerStr, ":")
									if len(resStr) == 2 {
										if strings.ToLower(resStr[0]) == "status" {
											statusRes := strings.Split(resStr[1], " ")

											statusCode, err := strconv.Atoi(statusRes[0])

											if err != nil {
												statusCode = 400
												log.Println("error when split status code")
											}

											fcgi.resp.StatusCode = statusCode
											fcgi.resp.Desc = statusRes[1]
										} else {
											fcgi.resp.AppendHeader(resStr[0], resStr[1])
										}

									}
									start = i + 1
									state = 1
								} else if state == 1 {
									stdoutHeader = true
									// 头部解析完毕
									if fcgi.resp.Body != nil {
										fcgi.resp.Body = append(fcgi.resp.Body, outdata[i+1:header.ContentLength]...)
									} else {
										fcgi.resp.Body = outdata[i+1 : header.ContentLength]
									}

									isFinish = true
								}
							}
							break

						default:
							state = 0
						}

						pre = v
						if isFinish {
							log.Println("解析完毕")
							break
						}
					}
				}

				break

			case FCGIStderr: // 错误输出流
				log.Printf("读取到错误流:%s\n", string(data[0:header.ContentLength]))
				break

			case FCGIEndRequest: // 结束请求
				log.Printf("protocal status: %v\n", data[4])
				stdoutHeader = false
				log.Println("释放阻塞")
				fcgi.finishChan <- true // 释放阻塞
				break
			}

			isHeader = true

			readLen = 8

		}

	}
}
