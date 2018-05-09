package http

import (
	"bytes"
	"io"
	"log"
)

type Request struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
	Body    []byte
}

// ParseHeader 函数用来解析输入流来构造请求头
func (req *Request) Parse(reader io.Reader) {

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	var in [1024]byte
	var buf bytes.Buffer
	var k string
	var v string
	var last byte

	firstline := true         //是否在匹配第一行
	endline := false          //是否匹配完一行
	kend := false             //当前行是否已经匹配到: 了
	info := make([]string, 0) //存储method url和version

	for {

		_, err := reader.Read(in[:])

		if err != nil {
			log.Fatal(err)
		}

		// 一个简单的有限状态机
		for index, c := range in {
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
							req.Body = in[index+1:] //构造body
							break                   // end for
						}
						endline = true
						kend = false
					} else {
						log.Fatal("unexpect \\n")
					}
					last = c
				}
			case ':':
				if !kend {
					k = buf.String()
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
