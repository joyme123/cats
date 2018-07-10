package http

import (
	"bytes"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joyme123/cats/config"
)

// Context , server的context，用来保存一些配置
type Context struct {
	KeyValue map[string]interface{}
}

// Handler 一次连接的处理句柄
type Handler struct {
	Response Response
	Request  Request
	conn     net.Conn
	srv      *Server
}

//Server 的封装
type Server struct {
	vhost   *config.VHost
	hub     Hub
	context Context
}

func (srv *Server) Context(vhost *config.VHost) {
	srv.vhost = vhost
	srv.context.KeyValue = make(map[string]interface{})
	srv.context.KeyValue["Addr"] = vhost.Addr
	srv.context.KeyValue["Port"] = strconv.Itoa(vhost.Port)
}

func (srv *Server) GetContext() *Context {
	return &(srv.context)
}

// 向server中注入组件
func (srv *Server) Register(component interface{}) {
	srv.hub.Register(component.(Component))
}

// Start the server
func (srv *Server) Start() {

	count := 0

	listener, err := net.Listen("tcp", srv.context.KeyValue["Addr"].(string)+":"+srv.context.KeyValue["Port"].(string))

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		count++

		var handler Handler

		handler.conn = conn
		handler.Response.Writer = conn
		handler.srv = srv

		go handler.Parse()
	}

}

// Parse 函数用来解析输入流来构造请求头
func (srv *Handler) Parse() {

	req := &(srv.Request)

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

		n, err := srv.conn.Read(in[:])

		if err != nil {
			log.Println(err)
			return
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
						log.Println("request parse error")
						return
					}
					req.Method = info[0]
					req.URI = info[1]
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
								// TODO:这个地方应该调用response进行输出
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
						log.Println("unexpect endline")
						return
					}
					last = c
				}
			case ':':
				if firstline {
					// 匹配第一行
					buf.WriteByte(c)
				} else {
					if !kend {
						k = strings.ToLower(buf.String())
						buf.Reset()
						kend = true
					} else {
						// 当前行已经匹配到:,说明当前匹配在value中
						buf.WriteByte(c)
					}
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
	}
}

func (srv *Handler) close() {
	// 检查该http请求的版本
	if srv.Request.Version == "HTTP/1.0" {
		if v, ok := srv.Request.Headers["connection"]; ok {
			if strings.ToLower(v) != "keep-alive" {
				// 不是长连接,断开
				srv.conn.Close()
			}
		} else {
			// 不是长连接,断开
			srv.conn.Close()
		}
	} else {
		if v, ok := srv.Request.Headers["connection"]; ok {
			if strings.ToLower(v) == "close" {
				// 不是长连接,断开
				srv.conn.Close()
			}
		}
	}
}

func (handler *Handler) Process() {

	defer handler.close()

	// 解析完毕一个请求
	if config.GetInstance().Log {
		handler.Request.logger(os.Stdout)
	}

	handler.Response.Version = handler.Request.Version

	for _, comp := range handler.srv.hub.container {
		comp.Serve(&handler.Request, &handler.Response)
	}

	handler.Response.out()
}
