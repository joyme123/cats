package http

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Handler struct {
	Response Response
	Request  Request
}

//Server 的封装
type Server struct {
}

type Config struct {
}

func (srv *Server) Config(config Config) {

}

// Start the server
func (srv *Server) Start() {

	listener, err := net.Listen("tcp", "127.0.0.1:8089")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		var handler Handler

		handler.Response.Writer = conn

		go handler.Parse(conn)
	}

}

// Parse 函数用来解析输入流来构造请求头
func (srv *Handler) Parse(reader io.Reader) {

	req := srv.Request

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	var in [1024]byte
	var buf bytes.Buffer
	var k string // header 的name,不区分大小写,这里统一转换为小写
	var v string // header 的值,区分大小写
	var last byte

	bodyLen := 0
	offset := 0

	firstline := true         //是否在匹配第一行
	endline := false          //是否匹配完一行
	kend := false             //当前行是否已经匹配到: 了
	info := make([]string, 0) //存储method url和version

	for {

		n, err := reader.Read(in[:])

		if err != nil {
			log.Fatal(err)
		}

		if bodyLen > 0 {
			if bodyLen > n {
				req.Body = append(req.Body, in[0:n-1]...)
				bodyLen -= n
				offset = n
			} else {
				req.Body = append(req.Body, in[0:bodyLen]...) // body解析结束
				bodyLen = 0
				offset = bodyLen
				firstline = true //是否在匹配第一行
				endline = false  //是否匹配完一行
				kend = false     //当前行是否已经匹配到: 了
				info = make([]string, 0)

				srv.Process()
			}

		}

		// 一个简单的有限状态机
		for index := offset; index < n; index++ {
			c := in[index]
			switch c {
			case '\r':
				if firstline {
					if len(info) < 2 {
						log.Fatal("request parse error")
					}
					req.Method = info[0]
					req.URL = info[1]
					req.Version = buf.String()
					buf.Reset()
				} else {

					if last != '\n' { //不是连续两次\r\n
						v = buf.String()
						buf.Reset()
						req.Headers[k] = v
					}
				}

				last = c
			case '\n':
				if firstline {
					firstline = false
				} else {
					if last == '\r' { //上一次是\r,则代表一行的结束
						if endline { // 处于endline状态,header解析完毕

							// 找到header中的Content-Length
							// log.Printf("header %v", req.Headers)
							bodyLen, err = strconv.Atoi(req.Headers["content-length"])
							if err != nil {
								// TODO 这个地方应该调用response进行输出
								// log.Fatal("err content length")
								bodyLen = 0
							}

							if bodyLen > 0 {

								if bodyLen+index <= n {
									req.Body = append(req.Body, in[index+1:index+bodyLen]...) //构造body结束
									bodyLen = 0

									// 整个请求已经解析结束，调用Process去处理
									srv.Process()
								} else {
									req.Body = append(req.Body, in[index+1:n]...) //构造body还没结束,但是in中的输入已经结束了
									bodyLen = bodyLen - (n - index - 1)
								}
							} else {
								// 整个请求已经解析结束，调用Process去处理
								srv.Process()
							}
							firstline = true //是否在匹配第一行
							endline = false  //是否匹配完一行
							kend = false     //当前行是否已经匹配到: 了
							info = make([]string, 0)
						} else {
							endline = true
							kend = false
						}
					} else {
						log.Fatal("unexpect \\n")
					}
					last = c
				}
			case ':':
				if !kend {
					k = strings.ToLower(buf.String())
					buf.Reset()
					kend = true
				}
				last = c
			case ' ':
				if firstline {
					info = append(info, buf.String())
					buf.Reset()
				} else {
					last = c
					endline = false
				}
			default:
				buf.WriteByte(c)
				last = c
				endline = false
			}
		} // end for state machine

		// 解析完毕一个请求
	}
}

func (srv *Handler) Process() {
	srv.Response.out()
}
