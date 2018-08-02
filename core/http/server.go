package http

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/joyme123/cats/config"
)

// Context , server的context，用来保存一些配置
// 全局使用的有:
// FilePath代表指向的文件路径
// IndexFiles代表索引的文件
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
//vhost是当前server的vhost配置
//hub中保存的是当前server中注入的组件
type Server struct {
	Addr  string
	Port  int
	sites []*config.Site
	vhs   map[string]*VirtualHost
}

func (srv *Server) Config(addr string, port int) {
	srv.Addr = addr
	srv.Port = port
}

func (srv *Server) AddSite(site *config.Site) {
	srv.sites = append(srv.sites, site)
}

func (srv *Server) Init() {
	srv.vhs = make(map[string]*VirtualHost)
}

func (srv *Server) GetVirtualHost() map[string]*VirtualHost {
	return srv.vhs
}

func (srv *Server) GetSite() []*config.Site {
	return srv.sites
}

func (srv *Server) SetVirtualHost(serverName string, vh *VirtualHost) {
	srv.vhs[serverName] = vh
}

// Start the server
func (srv *Server) Start() {

	listener, err := net.Listen("tcp", srv.Addr+":"+strconv.Itoa(srv.Port))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server start on %s:%d\n", srv.Addr, srv.Port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		var handler Handler

		// 初始化响应对象
		handler.Init(conn, srv)

		go handler.Parse()
	}

}

/**
 * 根据host去找对应的virtualhost
 *
 */
func (srv *Server) findHost(host string) *VirtualHost {

	// host要去掉端口号
	arr := strings.Split(host, ":")

	if len(arr) == 2 {
		host = arr[0]
	}

	log.Println("匹配值是" + host)
	if vh, ok := srv.vhs[host]; ok {
		return vh
	} else {
		return srv.vhs[srv.sites[0].ServerName]
	}
}

func (handler *Handler) Init(conn net.Conn, srv *Server) {
	handler.Response.Init(conn)
	handler.conn = conn
	handler.srv = srv
}

// Parse 函数用来解析输入流来构造请求头
func (handler *Handler) Parse() {
	defer handler.close()

	req := &(handler.Request)

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	if req.Context == nil {
		req.Context = make(map[string]interface{})
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

		n, err := handler.conn.Read(in[:])

		if err != nil {
			log.Println(err)
			return
		}

		parseFinish := false

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

				parseFinish = true
				handler.serverVh()
			}

		}

		// 一个简单的有限状态机
		for index := offset; index < n; index++ {
			if parseFinish {
				// body已经解析结束了，直接break出去
				break
			}
			c := in[index]
			switch c {
			case '\r':
				if firstline {
					if len(info) < 2 {
						log.Println("request parse error")
						return
					}
					req.Method = info[0]
					// 如果Method是Get，需要切割出QueryString
					if req.Method == "GET" {
						uriAndQuery := strings.SplitN(info[1], "?", 2) // 最多切割两个子串
						if len(uriAndQuery) == 2 {
							req.URI = uriAndQuery[0]
							req.QueryString = uriAndQuery[1]
						} else {
							req.URI = uriAndQuery[0]
						}

					} else {
						req.URI = info[1]
					}
					req.Version = buf.String()

					info = make([]string, 0) // 重置info
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

									parseFinish = true
									handler.serverVh()
								} else {
									req.Body = append(req.Body, in[index+1:n]...) //构造body还没结束,但是in中的输入已经结束了
									bodyLen = bodyLen - (n - index - 1)
								}
							} else {
								handler.serverVh()
								parseFinish = true
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

func (handler *Handler) serverVh() {
	// 整个请求已经解析结束，调用Process去处理
	vh := handler.srv.findHost(handler.Request.Headers["host"])
	if vh != nil {
		vh.ServeHttp(&(handler.Request), &(handler.Response))
	} else {
		// TODO: 这里如果请求的vhost不存在，要接输出错误页面
		log.Println("无对应serverName的服务")
	}
}

func (handler *Handler) close() {
	// 检查该http请求的版本
	if handler.Request.Version == "HTTP/1.0" {
		if v, ok := handler.Request.Headers["connection"]; ok {
			if strings.ToLower(v) != "keep-alive" {
				// 不是长连接,断开
				handler.conn.Close()
			}
		} else {
			// 不是长连接,断开
			handler.conn.Close()
		}
	} else {
		if v, ok := handler.Request.Headers["connection"]; ok {
			if strings.ToLower(v) == "close" {
				// 不是长连接,断开
				handler.conn.Close()
			}
		}
	}
}
