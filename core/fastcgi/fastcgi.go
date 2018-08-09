package fastcgi

// 目前的实现中为了简洁没有复用连接

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

// FastCGI 的结构体
type FastCGI struct {
	Index      int
	sockAdrr   string             // fastcgi应用程序的socket地址
	Context    *http.VhostContext // handler的context
	RootDir    string             // 根目录地址
	serverName string
	serverAddr string
	serverPort int
}

func (fcgi *FastCGI) commonHeaders(resp *http.Response) {
	resp.AppendHeader("connection", "keep-alive")
	resp.AppendHeader("server", "cats")
}

// New 方法是FastCGI 的实例化
func (fcgi *FastCGI) New(site *config.Site, context *http.VhostContext) {
	fcgi.sockAdrr = site.FCGIPass
	fcgi.Context = context
	fcgi.RootDir = site.Root
	fcgi.serverName = site.ServerName
	fcgi.serverAddr = site.Addr
	fcgi.serverPort = site.Port
}

// Start 方法是FastCGI在服务启动时调用的方法
func (fcgi *FastCGI) Start() {
	conn, err := net.Dial("tcp", fcgi.sockAdrr)
	if err != nil {
		log.Println("error when connect to FastCGI Application", err.Error())
		return
	}

	// 启动时向FCGI应用程序进行查询
	var getValuesRecord FCGIGetValuesRecord
	getValues := make(map[string]string)
	getValues[FCGIMaxConns] = ""
	getValues[FCGIMaxReqs] = ""
	getValues[FCGIMpxConns] = ""

	getValuesRecord.New(getValues)
	fcgi.sendRecord(conn, &getValuesRecord)

	readChan := make(chan bool)

	go fcgi.readManageInfo(conn, readChan)

	<-readChan

	conn.Close()
}

// Serve 方法是FastCGI在有请求到来时被调用的方法
func (fcgi *FastCGI) Serve(req *http.Request, resp *http.Response) {

	if resp.StatusCode != 0 {
		return
	}

	finishChan := make(chan bool)

	// 建立连接
	conn, err := net.Dial("tcp", fcgi.sockAdrr)
	if err != nil {
		log.Println("error when connect to FastCGI Application", err.Error())
		conn = nil
		resp.Error502()
		return
	}

	log.Println("新建的fcgi连接")

	go fcgi.readHandler(conn, req, resp, finishChan)

	currentID := uint16(1)
	log.Printf("请求ID:%v\n", currentID)
	// 1.创建并发送开始请求
	var beginRequestRecord FCGIBeginRequestRecord
	beginRequestRecord.New(currentID, FCGIResponder)

	fcgi.sendRecord(conn, &beginRequestRecord)

	// 2.获取当前请求的请求头，将其传递给fastcgi 程序
	var paramsRecord FCGIParamsRecord
	var params map[string]string
	params = make(map[string]string)

	filepath, ok := req.Context["FilePath"]

	if !ok {
		log.Println("serve file error: not found filepath in request context")
		resp.Error404()
		return
	}

	if filepathStr, ok := filepath.(string); ok {

		params["SCRIPT_FILENAME"] = filepathStr
		params["QUERY_STRING"] = req.QueryString
		params["REQUEST_METHOD"] = req.Method
		params["SCRIPT_NAME"] = req.URI
		params["REQUEST_URI"] = req.URI
		params["DOCUMENT_URI"] = req.URI
		params["DOCUMENT_ROOT"] = fcgi.RootDir
		params["SERVER_PROTOCOL"] = req.Version
		params["GATEWAY_INTERFACE"] = "CGI/1.1"
		params["SERVER_SOFTWARE"] = "cats"

		remoteAddr := strings.Split(req.RemoteAddr, ":")

		params["REMOTE_ADDR"] = remoteAddr[0]
		params["REMOTE_PORT"] = remoteAddr[1]
		params["SERVER_ADDR"] = fcgi.serverAddr
		params["SERVER_PORT"] = strconv.Itoa(fcgi.serverPort)
		params["SERVER_NAME"] = fcgi.serverName

		if accept, ok := req.Headers["accept"]; ok {
			params["HTTP_ACCEPT"] = accept
		}

		if acceptLang, ok := req.Headers["accept-language"]; ok {
			params["HTTP_ACCEPT_LANGUAGE"] = acceptLang
		}

		if acceptEnc, ok := req.Headers["accept-encoding"]; ok {
			params["HTTP_ACCEPT_ENCODING"] = acceptEnc
		}

		if userAgent, ok := req.Headers["user-agent"]; ok {
			params["HTTP_USER_AGENT"] = userAgent
		}

		if host, ok := req.Headers["host"]; ok {
			params["HTTP_HOST"] = host
		}

		if connection, ok := req.Headers["connection"]; ok {
			params["HTTP_CONNECTION"] = connection
		}

		if contentType, ok := req.Headers["content-type"]; ok {
			params["HTTP_CONTENT_TYPE"] = contentType
			params["CONTENT_TYPE"] = contentType
		}

		if contentLength, ok := req.Headers["content-length"]; !ok {
			params["CONTENT_LENGTH"] = "0"
			params["HTTP_CONTENT_LENGTH"] = "0"
		} else {
			params["CONTENT_LENGTH"] = contentLength
			params["HTTP_CONTENT_LENGTH"] = contentLength
		}

		if cacheCtrl, ok := req.Headers["cache-control"]; ok {
			params["HTTP_CACHE_CONTROL"] = cacheCtrl
		}

		if cookie, ok := req.Headers["cookie"]; ok {
			params["HTTP_COOKIE"] = cookie
		}

		paramsRecord.New(currentID, params)
		fcgi.sendRecord(conn, &paramsRecord)

		var emptyParamsRecord FCGIParamsRecord
		emptyParams := make(map[string]string)
		emptyParamsRecord.New(currentID, emptyParams)
		fcgi.sendRecord(conn, &emptyParamsRecord)

		// 3.创建并发送stdin请求
		var stdinRecord FCGIStdinRecord
		log.Printf("请求体的内容:%s", req.Body)
		stdinRecord.New(currentID, req.Body)
		fcgi.sendRecord(conn, &stdinRecord)

		if len(req.Body) != 0 {
			var emptyStdinRecord FCGIStdinRecord
			emptyBytes := make([]byte, 0)
			emptyStdinRecord.New(currentID, emptyBytes)
			fcgi.sendRecord(conn, &emptyStdinRecord)

			fcgi.sendRecord(conn, &emptyStdinRecord)
		}

		// 这里应该阻塞起来等待fastcgi程序响应

		<-finishChan // 使用管道阻塞
		conn.Close()
		log.Println("阻塞结束")
		fcgi.commonHeaders(resp)

		if resp.StatusCode == 0 {
			resp.StatusCode = 200
			resp.Desc = "OK"
		}
	} else {
		resp.Error404()
	}

}

// Shutdown 方法是FastCGI在服务终止时被调用的方法
func (fcgi *FastCGI) Shutdown() {
	log.Println("关闭fastcgi连接")
}

// GetIndex 用来获取当前组件的索引
func (fcgi *FastCGI) GetIndex() int {
	return fcgi.Index
}

func (fcgi *FastCGI) GetContainer() string {
	return "location"
}

// sendRecord 发送FastCGI 记录
func (fcgi *FastCGI) sendRecord(conn net.Conn, record Record) {

	// 通过fcgiConn发送记录
	_, err := conn.Write(record.ToBlob())

	if err != nil {
		log.Printf("fcgi 写入错误: %v\n", err)
	}

}

func (fcgi *FastCGI) readManageInfo(conn net.Conn, readChan chan bool) {
	var header FCGIHeader
	readLen := 8
	isHeader := true

	for {
		data := make([]byte, readLen, readLen)
		n, err := conn.Read(data[:])

		if err != nil {
			log.Println("error when read fcgi manage info", err.Error())

			return
		}

		if n == 0 {
			log.Println("没有读取数据")
			readChan <- true
			return
		}

		if isHeader {
			header.New(data[0:readLen]) // 初始化header
			isHeader = false
			readLen = int(header.ContentLength) + int(header.PaddingLength)
		} else {
			if header.Type == FCGIGetValueResults {

				pair := fcgi.parseNameValuePair(data, int(header.ContentLength))

				for key, value := range pair {
					switch key {
					case FCGIMaxConns: // 最大连接数
						log.Printf("最大连接数:%v\n", value)
						break

					case FCGIMaxReqs: // 最多请求数
						log.Printf("最多请求数:%v\n", value)
						break

					case FCGIMpxConns: // 是否复用连接,0不复用连接
						log.Printf("是否复用连接:%v\n", value)
						break
					}
				}
			} else {
				log.Println("fcgi应用程序回复错误")
			}

			readChan <- true
		}
	}
}

// readHandler 负责从FastCGI应用程序中读取stdout,stderr,以及EndRequestRecord
func (fcgi *FastCGI) readHandler(conn net.Conn, req *http.Request, resp *http.Response, finishChan chan bool) {
	var header FCGIHeader
	readLen := 8 // 下一次要读取的长度
	isHeader := true
	stdoutHeader := false // 是否解析过标准输出的头信息
	for {
		data := make([]byte, readLen, readLen)
		n, err := conn.Read(data[:])

		if err != nil {
			log.Println("error when read data from FastCGI Application", err.Error())
			finishChan <- true // 释放阻塞
			return
		}

		if n == 0 {
			log.Println("没有读取数据")
			finishChan <- true // 释放阻塞
			return
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
					if resp.Body != nil {
						resp.Body = append(resp.Body, outdata[:header.ContentLength]...)
					} else {
						resp.Body = outdata[:header.ContentLength]
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
									pos := strings.IndexByte(headerStr, ':')
									var headerKey string
									var headerValue string
									if pos > 0 {
										headerKey = headerStr[0:pos]
										headerValue = headerStr[pos+1:]
									}
									if headerKey != "" && headerValue != "" {
										if strings.ToLower(headerKey) == "status" {
											statusRes := strings.Split(strings.Trim(headerValue, " "), " ")

											statusCode, err := strconv.Atoi(statusRes[0])

											if err != nil {
												statusCode = 400
												log.Println("error when split status code")
											}

											resp.StatusCode = statusCode
											resp.Desc = statusRes[1]
										} else {
											resp.AppendHeader(strings.ToLower(headerKey), headerValue)
										}

										headerKey = ""
										headerValue = ""

									}
									start = i + 1
									state = 1
								} else if state == 1 {
									stdoutHeader = true
									// 头部解析完毕
									if resp.Body != nil {
										resp.Body = append(resp.Body, outdata[i+1:header.ContentLength]...)
									} else {
										resp.Body = outdata[i+1 : header.ContentLength]
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
				switch data[4] {
				case FCGIRequestComplete:
					// 正常的请求结束。
					break

				case FCGICantMpxConn:
					// 拒绝新请求。当Web服务器通过一个连接将并发请求发送到旨在每个连接一次处理一个请求的应用程序时，就会发生这种情况。
					break

				case FCGIOverloaded:
					// 拒绝新请求。当应用程序耗尽某些资源时会发生这种情况，例如：数据库连接。
					break

				case FCGIUnknownRole:
					// 拒绝新请求。当Web服务器指定了应用程序未知的角色时，会发生这种情况。
					break
				}
				stdoutHeader = false
				log.Println("释放阻塞")
				finishChan <- true // 释放阻塞
				break
			}

			isHeader = true

			readLen = 8

		}

	}
}

func (fcgi *FastCGI) parseNameValuePair(data []byte, len int) map[string]string {
	pair := make(map[string]string)
	if len <= 0 {
		return pair
	}

	offset := 0
	var keyLen uint32
	var valueLen uint32

	for offset < len {
		if data[offset]>>7 == 0 {
			keyLen = uint32(data[offset])
			// keylen一个字节
			offset = offset + 1

			if data[offset]>>7 == 0 {
				// valueLen 一个字节
				valueLen = uint32(data[offset])
				offset = offset + 1
			} else {
				//valueLen 4个字节
				data[offset] = data[offset] & 0x7f
				binary.BigEndian.PutUint32(data[offset:offset+4], valueLen)
				offset = offset + 4
			}

			key := data[offset : offset+int(keyLen)]
			offset += int(keyLen)

			value := data[offset : offset+int(valueLen)]
			offset += int(valueLen)

			pair[string(key)] = string(value)
		} else {
			// keyLen4个字节
			data[offset] = data[offset] & 0x7f
			binary.BigEndian.PutUint32(data[offset:offset+4], keyLen)
			offset = offset + 4
			if data[offset]>>7 == 0 {
				// valueLen 一个字节
				valueLen = uint32(data[offset])
				offset++
			} else {
				// valueLen 4个字节
				data[offset] = data[offset] & 0x7f
				binary.BigEndian.PutUint32(data[offset:offset+4], valueLen)
				offset = offset + 4
			}

			key := data[offset : offset+int(keyLen)]
			offset += int(keyLen)

			value := data[offset : offset+int(valueLen)]
			offset += int(valueLen)
			pair[string(key)] = string(value)
		}
	}

	return pair
}
